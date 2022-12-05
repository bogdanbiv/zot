// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql_generated

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Annotation struct {
	Key   *string `json:"Key"`
	Value *string `json:"Value"`
}

// Contains various details about the CVE and a list of PackageInfo about the affected packages
type Cve struct {
	ID          *string        `json:"Id"`
	Title       *string        `json:"Title"`
	Description *string        `json:"Description"`
	Severity    *string        `json:"Severity"`
	PackageList []*PackageInfo `json:"PackageList"`
}

// Contains the tag of the image and a list of CVEs
type CVEResultForImage struct {
	Tag     *string   `json:"Tag"`
	CVEList []*Cve    `json:"CVEList"`
	Page    *PageInfo `json:"Page"`
}

// Parameters used to filter the results of the query
type Filter struct {
	Os            []*string `json:"Os"`
	Arch          []*string `json:"Arch"`
	HasToBeSigned *bool     `json:"HasToBeSigned"`
}

// Search everything. Can search Images, Repos and Layers
type GlobalSearchResult struct {
	Page   *PageInfo       `json:"Page"`
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

// Contains details about the image
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
	Authors         *string                    `json:"Authors"`
	Logo            *string                    `json:"Logo"`
}

type ImageVulnerabilitySummary struct {
	MaxSeverity *string `json:"MaxSeverity"`
	Count       *int    `json:"Count"`
}

type LayerHistory struct {
	Layer              *LayerSummary       `json:"Layer"`
	HistoryDescription *HistoryDescription `json:"HistoryDescription"`
}

// Contains details about the layer
type LayerSummary struct {
	Size   *string `json:"Size"`
	Digest *string `json:"Digest"`
	Score  *int    `json:"Score"`
}

// Contains details about the supported OS and architecture of the image
type OsArch struct {
	Os   *string `json:"Os"`
	Arch *string `json:"Arch"`
}

// Contains the name of the package, the current installed version and the version where the CVE was fixed
type PackageInfo struct {
	Name             *string `json:"Name"`
	InstalledVersion *string `json:"InstalledVersion"`
	FixedVersion     *string `json:"FixedVersion"`
}

type PageInfo struct {
	TotalCount int `json:"TotalCount"`
	ItemCount  int `json:"ItemCount"`
}

// Pagination parameters.
// Limit: refers to the amout of results per page. If you set limit to -1, the pagination behaior is disabled
// Offset: the results page number you want to receive.
// Sort by: the criteria used to sort the results on the page.
type PageInput struct {
	Limit  *int          `json:"limit"`
	Offset *int          `json:"offset"`
	SortBy *SortCriteria `json:"sortBy"`
}

type PaginatedImagesResult struct {
	Page    *PageInfo       `json:"Page"`
	Results []*ImageSummary `json:"Results"`
}

type PaginatedReposResult struct {
	Page    *PageInfo      `json:"Page"`
	Results []*RepoSummary `json:"Results"`
}

type Referrer struct {
	MediaType    *string       `json:"MediaType"`
	ArtifactType *string       `json:"ArtifactType"`
	Size         *int          `json:"Size"`
	Digest       *string       `json:"Digest"`
	Annotations  []*Annotation `json:"Annotations"`
}

// Contains details about the repo: a list of image summaries and a summary of the repo
type RepoInfo struct {
	Images  []*ImageSummary `json:"Images"`
	Summary *RepoSummary    `json:"Summary"`
}

// Contains details about the repo
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
	IsStarred     *bool         `json:"IsStarred"`
}

type SortCriteria string

const (
	SortCriteriaRelevance     SortCriteria = "RELEVANCE"
	SortCriteriaUpdateTime    SortCriteria = "UPDATE_TIME"
	SortCriteriaAlphabeticAsc SortCriteria = "ALPHABETIC_ASC"
	SortCriteriaAlphabeticDsc SortCriteria = "ALPHABETIC_DSC"
	SortCriteriaSeverity      SortCriteria = "SEVERITY"
	SortCriteriaStars         SortCriteria = "STARS"
	SortCriteriaDownloads     SortCriteria = "DOWNLOADS"
)

var AllSortCriteria = []SortCriteria{
	SortCriteriaRelevance,
	SortCriteriaUpdateTime,
	SortCriteriaAlphabeticAsc,
	SortCriteriaAlphabeticDsc,
	SortCriteriaSeverity,
	SortCriteriaStars,
	SortCriteriaDownloads,
}

func (e SortCriteria) IsValid() bool {
	switch e {
	case SortCriteriaRelevance, SortCriteriaUpdateTime, SortCriteriaAlphabeticAsc, SortCriteriaAlphabeticDsc, SortCriteriaSeverity, SortCriteriaStars, SortCriteriaDownloads:
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
