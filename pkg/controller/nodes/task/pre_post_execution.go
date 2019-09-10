package task

import (
	"context"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	pluginCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/io"
	"github.com/lyft/flytestdlib/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lyft/flytepropeller/pkg/controller/nodes/task/catalog"
)

func (t *Handler) CheckCatalogCache(ctx context.Context, tr pluginCore.TaskReader, inputReader io.InputReader, outputWriter io.OutputWriter) (bool, error) {
	tk, err := tr.Read(ctx)
	if err != nil {
		logger.Errorf(ctx, "Failed to read TaskTemplate, error :%s", err.Error())
		return false, err
	}
	if tk.Metadata.Discoverable {
		key := catalog.Key{
			Identifier:     *tk.Id,
			CacheVersion:   tk.Metadata.DiscoveryVersion,
			TypedInterface: *tk.Interface,
			InputReader:    inputReader,
		}
		if resp, err := t.catalog.Get(ctx, key); err != nil {
			if taskStatus, ok := status.FromError(err); ok && taskStatus.Code() == codes.NotFound {
				t.metrics.discoveryMissCount.Inc(ctx)
				logger.Infof(ctx, "Artifact not found in Discovery. Executing Task.")
			} else {
				t.metrics.discoveryGetFailureCount.Inc(ctx)
				logger.Errorf(ctx, "Discovery check failed. Executing Task. Err: %v", err.Error())
			}
		} else if resp != nil {
			t.metrics.discoveryHitCount.Inc(ctx)
			if iface := tk.Interface; iface != nil && iface.Outputs != nil && len(iface.Outputs.Variables) > 0 {
				if err := outputWriter.Put(ctx, resp); err != nil {
					logger.Errorf(ctx, "failed to write data to Storage, err: %v", err.Error())
					return false, errors.Wrapf(err, "failed to copy cached results for task.")
				}
			}
			// SetCached.
			return true, nil
		} else {
			// Nil response and Nil error
			t.metrics.discoveryGetFailureCount.Inc(ctx)
			return false, errors.Wrapf(err, "Nil catalog response. Failed to check Catalog for previous results")
		}
	}
	return false, nil
}

func (t *Handler) ValidateOutputAndCacheAdd(ctx context.Context, i io.InputReader, r io.OutputReader, tr pluginCore.TaskReader, m catalog.Metadata) (*io.ExecutionError, error) {
	// access w
	y, err := r.IsError(ctx)
	if err != nil {
		return nil, err
	}
	if y {
		taskErr, err := r.ReadError(ctx)
		if err != nil {
			return nil, err
		}
		return &taskErr, nil
	}

	tk, err := tr.Read(ctx)
	if err != nil {
		logger.Errorf(ctx, "Failed to read TaskTemplate, error :%s", err.Error())
		return nil, err
	}
	// Do this if we have outputs declared for the Handler interface!
	if iface := tk.Interface; iface != nil && iface.Outputs != nil && len(iface.Outputs.Variables) > 0 {
		ok, err := r.Exists(ctx)
		if err != nil {
			logger.Errorf(ctx, "Failed to check if the output file exists. Error: %s", err.Error())
			return nil, err
		}
		if !ok {
			// Does not exist
			return &io.ExecutionError{
				ExecutionError: &core.ExecutionError{
					Code:    "OutputsNotFound",
					Message: "Outputs not generated by task execution",
				},
				IsRecoverable: true,
			}, nil
		}

		if !r.IsFile(ctx) {
			// Read output and write to file
			logger.Warnf(ctx, "Inputs of type file are only handled currently. Implement other input types")
		}

		// ignores discovery write failures
		if tk.Metadata.Discoverable {
			key := catalog.Key{
				Identifier:     *tk.Id,
				CacheVersion:   "",
				TypedInterface: *tk.Interface,
				InputReader:    i,
			}
			if err2 := t.catalog.Put(ctx, key, r, m); err2 != nil {
				t.metrics.discoveryPutFailureCount.Inc(ctx)
				logger.Errorf(ctx, "Failed to write results to catalog. Err: %v", err2)
			} else {
				logger.Debugf(ctx, "Successfully cached results to discovery - Task [%s]", tk.GetId())
			}
		}
	}
	return nil, nil
}
