// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql_generated

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Paginated interface {
	IsPaginated()
	GetPage() *PageInfo
}

type Cve struct {
	ID          *string        `json:"Id"`
	Title       *string        `json:"Title"`
	Description *string        `json:"Description"`
	Severity    *string        `json:"Severity"`
	PackageList []*PackageInfo `json:"PackageList"`
}

type CVEResultForImage struct {
	Tag     *string `json:"Tag"`
	CVEList []*Cve  `json:"CVEList"`
}

type GlobalSearchResult struct {
	Images []*ImageSummary `json:"Images"`
	Repos  []*RepoSummary  `json:"Repos"`
	Layers []*LayerSummary `json:"Layers"`
}

type HistoryDescription struct {
	Created *time.Time `json:"Created"`
	// CreatedBy is the command which created the layer.
	CreatedBy *string `json:"CreatedBy"`
	// Author is the author of the build point.
	Author *string `json:"Author"`
	// Comment is a custom message set when creating the layer.
	Comment *string `json:"Comment"`
	// EmptyLayer is used to mark if the history item created a filesystem diff.
	EmptyLayer *bool `json:"EmptyLayer"`
}

type ImageSummary struct {
	RepoName        *string                    `json:"RepoName"`
	Tag             *string                    `json:"Tag"`
	Digest          *string                    `json:"Digest"`
	ConfigDigest    *string                    `json:"ConfigDigest"`
	LastUpdated     *time.Time                 `json:"LastUpdated"`
	IsSigned        *bool                      `json:"IsSigned"`
	Size            *string                    `json:"Size"`
	Platform        *OsArch                    `json:"Platform"`
	Vendor          *string                    `json:"Vendor"`
	Score           *int                       `json:"Score"`
	DownloadCount   *int                       `json:"DownloadCount"`
	Layers          []*LayerSummary            `json:"Layers"`
	Description     *string                    `json:"Description"`
	Licenses        *string                    `json:"Licenses"`
	Labels          *string                    `json:"Labels"`
	Title           *string                    `json:"Title"`
	Source          *string                    `json:"Source"`
	Documentation   *string                    `json:"Documentation"`
	History         []*LayerHistory            `json:"History"`
	Vulnerabilities *ImageVulnerabilitySummary `json:"Vulnerabilities"`
}

type ImageVulnerabilitySummary struct {
	MaxSeverity *string `json:"MaxSeverity"`
	Count       *int    `json:"Count"`
}

type LayerHistory struct {
	Layer              *LayerSummary       `json:"Layer"`
	HistoryDescription *HistoryDescription `json:"HistoryDescription"`
}

type LayerSummary struct {
	Size   *string `json:"Size"`
	Digest *string `json:"Digest"`
	Score  *int    `json:"Score"`
}

type MutationResult struct {
	// outcome of the Mutation
	Success bool `json:"success"`
}

type OsArch struct {
	Os   *string `json:"Os"`
	Arch *string `json:"Arch"`
}

type PackageInfo struct {
	Name             *string `json:"Name"`
	InstalledVersion *string `json:"InstalledVersion"`
	FixedVersion     *string `json:"FixedVersion"`
}

type PageInfo struct {
	ObjectCount  int  `json:"ObjectCount"`
	PreviousPage *int `json:"PreviousPage"`
	NextPage     *int `json:"NextPage"`
	Pages        *int `json:"Pages"`
}

type PageInput struct {
	Limit  *int          `json:"limit"`
	Offset *int          `json:"offset"`
	SortBy *SortCriteria `json:"sortBy"`
}

type PaginatedImagesResult struct {
	Page    *PageInfo       `json:"Page"`
	Results []*ImageSummary `json:"Results"`
}

func (PaginatedImagesResult) IsPaginated()            {}
func (this PaginatedImagesResult) GetPage() *PageInfo { return this.Page }

type PaginatedReposResult struct {
	Page    *PageInfo      `json:"Page"`
	Results []*RepoSummary `json:"Results"`
}

func (PaginatedReposResult) IsPaginated()            {}
func (this PaginatedReposResult) GetPage() *PageInfo { return this.Page }

type RepoInfo struct {
	Images  []*ImageSummary `json:"Images"`
	Summary *RepoSummary    `json:"Summary"`
}

type RepoSummary struct {
	Name          *string       `json:"Name"`
	LastUpdated   *time.Time    `json:"LastUpdated"`
	Size          *string       `json:"Size"`
	Platforms     []*OsArch     `json:"Platforms"`
	Vendors       []*string     `json:"Vendors"`
	Score         *int          `json:"Score"`
	NewestImage   *ImageSummary `json:"NewestImage"`
	DownloadCount *int          `json:"DownloadCount"`
	StarCount     *int          `json:"StarCount"`
	IsBookmarked  *bool         `json:"IsBookmarked"`
}

type SortCriteria string

const (
	SortCriteriaRelevance     SortCriteria = "RELEVANCE"
	SortCriteriaUpdateTime    SortCriteria = "UPDATE_TIME"
	SortCriteriaAlphabeticAsc SortCriteria = "ALPHABETIC_ASC"
	SortCriteriaAlphabeticDsc SortCriteria = "ALPHABETIC_DSC"
	SortCriteriaStars         SortCriteria = "STARS"
	SortCriteriaDownloads     SortCriteria = "DOWNLOADS"
)

var AllSortCriteria = []SortCriteria{
	SortCriteriaRelevance,
	SortCriteriaUpdateTime,
	SortCriteriaAlphabeticAsc,
	SortCriteriaAlphabeticDsc,
	SortCriteriaStars,
	SortCriteriaDownloads,
}

func (e SortCriteria) IsValid() bool {
	switch e {
	case SortCriteriaRelevance, SortCriteriaUpdateTime, SortCriteriaAlphabeticAsc, SortCriteriaAlphabeticDsc, SortCriteriaStars, SortCriteriaDownloads:
		return true
	}
	return false
}

func (e SortCriteria) String() string {
	return string(e)
}

func (e *SortCriteria) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortCriteria(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortCriteria", str)
	}
	return nil
}

func (e SortCriteria) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
