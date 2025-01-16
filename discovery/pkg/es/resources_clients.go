// Code is generated by go generate. DO NOT EDIT.
package opengovernance

import (
	"context"
	template "github.com/opengovern/og-describer-template/discovery/provider"
	essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"runtime"
)

type Client struct {
	essdk.Client
}

// ==========================  START: ArtifactDockerFile =============================

type ArtifactDockerFile struct {
	ResourceID      string                                 `json:"resource_id"`
	PlatformID      string                                 `json:"platform_id"`
	Description     template.ArtifactDockerFileDescription `json:"Description"`
	Metadata        template.Metadata                      `json:"metadata"`
	DescribedBy     string                                 `json:"described_by"`
	ResourceType    string                                 `json:"resource_type"`
	IntegrationType string                                 `json:"integration_type"`
	IntegrationID   string                                 `json:"integration_id"`
}

type ArtifactDockerFileHit struct {
	ID      string             `json:"_id"`
	Score   float64            `json:"_score"`
	Index   string             `json:"_index"`
	Type    string             `json:"_type"`
	Version int64              `json:"_version,omitempty"`
	Source  ArtifactDockerFile `json:"_source"`
	Sort    []interface{}      `json:"sort"`
}

type ArtifactDockerFileHits struct {
	Total essdk.SearchTotal       `json:"total"`
	Hits  []ArtifactDockerFileHit `json:"hits"`
}

type ArtifactDockerFileSearchResponse struct {
	PitID string                 `json:"pit_id"`
	Hits  ArtifactDockerFileHits `json:"hits"`
}

type ArtifactDockerFilePaginator struct {
	paginator *essdk.BaseESPaginator
}

func (k Client) NewArtifactDockerFilePaginator(filters []essdk.BoolFilter, limit *int64) (ArtifactDockerFilePaginator, error) {
	paginator, err := essdk.NewPaginator(k.ES(), "github_artifact_dockerfile", filters, limit)
	if err != nil {
		return ArtifactDockerFilePaginator{}, err
	}

	p := ArtifactDockerFilePaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p ArtifactDockerFilePaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p ArtifactDockerFilePaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p ArtifactDockerFilePaginator) NextPage(ctx context.Context) ([]ArtifactDockerFile, error) {
	var response ArtifactDockerFileSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []ArtifactDockerFile
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var listArtifactDockerFileFilters = map[string]string{
	"dockerfile_content": "Description.DockerfileContent",
	"html_url":           "Description.HTMLURL",
	"images":             "Description.Images",
	"last_updated_at":    "Description.LastUpdatedAt",
	"name":               "Description.Name",
	"repository":         "Description.Repository",
	"sha":                "Description.Sha",
}

func ListArtifactDockerFile(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ListArtifactDockerFile")
	runtime.GC()

	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactDockerFile NewClientCached", "error", err)
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactDockerFile NewSelfClientCached", "error", err)
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactDockerFile GetConfigTableValueOrNil for OpenGovernanceConfigKeyIntegrationID", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactDockerFile GetConfigTableValueOrNil for OpenGovernanceConfigKeyResourceCollectionFilters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactDockerFile GetConfigTableValueOrNil for OpenGovernanceConfigKeyClientType", "error", err)
		return nil, err
	}

	paginator, err := k.NewArtifactDockerFilePaginator(essdk.BuildFilter(ctx, d.QueryContext, listArtifactDockerFileFilters, integrationId, encodedResourceCollectionFilters, clientType), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactDockerFile NewArtifactDockerFilePaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListArtifactDockerFile paginator.NextPage", "error", err)
			return nil, err
		}

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getArtifactDockerFileFilters = map[string]string{
	"dockerfile_content": "Description.DockerfileContent",
	"html_url":           "Description.HTMLURL",
	"images":             "Description.Images",
	"last_updated_at":    "Description.LastUpdatedAt",
	"name":               "Description.Name",
	"repository":         "Description.Repository",
	"sha":                "Description.Sha",
}

func GetArtifactDockerFile(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetArtifactDockerFile")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		return nil, err
	}

	limit := int64(1)
	paginator, err := k.NewArtifactDockerFilePaginator(essdk.BuildFilter(ctx, d.QueryContext, getArtifactDockerFileFilters, integrationId, encodedResourceCollectionFilters, clientType), &limit)
	if err != nil {
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			return v, nil
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ==========================  END: ArtifactDockerFile =============================
