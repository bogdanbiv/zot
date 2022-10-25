package search

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"zotregistry.io/zot/pkg/extensions/search/common"
	"zotregistry.io/zot/pkg/extensions/search/gql_generated"
)

// CVEListForImage is the resolver for the CVEListForImage field.
func (r *queryResolver) CVEListForImage(ctx context.Context, image string) (*gql_generated.CVEResultForImage, error) {
	return getCVEListForImage(ctx, image, r.cveInfo, r.log)
}

// ImageListForCve is the resolver for the ImageListForCVE field.
func (r *queryResolver) ImageListForCve(ctx context.Context, id string, requestedPage *gql_generated.PageInput) ([]*gql_generated.ImageSummary, error) {
	return getImageListForCVE(ctx, id, r.cveInfo, requestedPage, r.repoDB, r.log)
}

// ImageListWithCVEFixed is the resolver for the ImageListWithCVEFixed field.
func (r *queryResolver) ImageListWithCVEFixed(ctx context.Context, id string, image string, requestedPage *gql_generated.PageInput) ([]*gql_generated.ImageSummary, error) {
	return getImageListWithCVEFixed(ctx, id, image, r.cveInfo, requestedPage, r.repoDB, r.log)
}

// ImageListForDigest is the resolver for the ImageListForDigest field.
func (r *queryResolver) ImageListForDigest(ctx context.Context, id string, requestedPage *gql_generated.PageInput) ([]*gql_generated.ImageSummary, error) {
	r.log.Info().Msg("extracting repositories")

	imgResultForDigest, err := getImageListForDigest(ctx, id, r.repoDB, r.cveInfo, requestedPage)

	return imgResultForDigest, err
}

// RepoListWithNewestImage is the resolver for the RepoListWithNewestImage field.
func (r *queryResolver) RepoListWithNewestImage(ctx context.Context, requestedPage *gql_generated.PageInput) ([]*gql_generated.RepoSummary, error) {
	r.log.Info().Msg("extension api: finding image list")

	reposSummary, err := repoListWithNewestImage(ctx, r.cveInfo, r.log, requestedPage, r.repoDB)
	if err != nil {
		r.log.Error().Err(err).Msg("unable to retrieve repo list")

		return reposSummary, err
	}

	return reposSummary, nil
}

// ImageList is the resolver for the ImageList field.
func (r *queryResolver) ImageList(ctx context.Context, repo string) ([]*gql_generated.ImageSummary, error) {
	r.log.Info().Msg("extension api: getting a list of all images")

	imageList := make([]*gql_generated.ImageSummary, 0)

	defaultStore := r.storeController.DefaultStore

	dsImageList, err := r.getImageList(defaultStore, repo)
	if err != nil {
		r.log.Error().Err(err).Msg("extension api: error extracting default store image list")

		return imageList, err
	}

	if len(dsImageList) != 0 {
		imageList = append(imageList, dsImageList...)
	}

	subStore := r.storeController.SubStore

	for _, store := range subStore {
		ssImageList, err := r.getImageList(store, repo)
		if err != nil {
			r.log.Error().Err(err).Msg("extension api: error extracting substore image list")

			return imageList, err
		}

		if len(ssImageList) != 0 {
			imageList = append(imageList, ssImageList...)
		}
	}

	return imageList, nil
}

// ExpandedRepoInfo is the resolver for the ExpandedRepoInfo field.
func (r *queryResolver) ExpandedRepoInfo(ctx context.Context, repo string) (*gql_generated.RepoInfo, error) {
	repoInfo, err := expandedRepoInfo(ctx, repo, r.repoDB, r.cveInfo, r.log)

	return repoInfo, err
}

// GlobalSearch is the resolver for the GlobalSearch field.
func (r *queryResolver) GlobalSearch(ctx context.Context, query string, filter *gql_generated.Filter, requestedPage *gql_generated.PageInput) (*gql_generated.GlobalSearchResult, error) {
	if err := validateGlobalSearchInput(query, filter, requestedPage); err != nil {
		return &gql_generated.GlobalSearchResult{}, err
	}

	query = cleanQuery(query)
	filter = cleanFilter(filter)

	repos, images, layers, err := globalSearch(ctx, query, r.repoDB, filter, requestedPage, r.cveInfo, r.log)

	return &gql_generated.GlobalSearchResult{
		Images: images,
		Repos:  repos,
		Layers: layers,
	}, err
}

// DependencyListForImage is the resolver for the DependencyListForImage field.
func (r *queryResolver) DerivedImageList(ctx context.Context, image string) ([]*gql_generated.ImageSummary, error) {
	layoutUtils := common.NewBaseOciLayoutUtils(r.storeController, r.log)
	imageList := make([]*gql_generated.ImageSummary, 0)

	repoList, err := layoutUtils.GetRepositories()
	if err != nil {
		r.log.Error().Err(err).Msg("unable to get repositories list")

		return nil, err
	}

	if len(repoList) == 0 {
		r.log.Info().Msg("no repositories found")

		return imageList, nil
	}

	imageDir, imageTag := common.GetImageDirAndTag(image)

	imageManifest, _, err := layoutUtils.GetImageManifest(imageDir, imageTag)
	if err != nil {
		r.log.Info().Str("image", image).Msg("image not found")

		return imageList, err
	}

	imageLayers := imageManifest.Layers

	for _, repo := range repoList {
		repoInfo, err := r.ExpandedRepoInfo(ctx, repo)
		if err != nil {
			r.log.Error().Err(err).Msg("unable to get image list")

			return nil, err
		}

		imageSummaries := repoInfo.Images

		// verify every image
		for _, imageSummary := range imageSummaries {
			if imageTag == *imageSummary.Tag && imageDir == repo {
				continue
			}

			layers := imageSummary.Layers

			sameLayer := 0

			for _, l := range imageLayers {
				for _, k := range layers {
					if *k.Digest == l.Digest.String() {
						sameLayer++
					}
				}
			}

			// if all layers are the same
			if sameLayer == len(imageLayers) {
				// add to returned list
				imageList = append(imageList, imageSummary)
			}
		}
	}

	return imageList, nil
}

// BaseImageList is the resolver for the BaseImageList field.
func (r *queryResolver) BaseImageList(ctx context.Context, image string) ([]*gql_generated.ImageSummary, error) {
	layoutUtils := common.NewBaseOciLayoutUtils(r.storeController, r.log)
	imageList := make([]*gql_generated.ImageSummary, 0)

	repoList, err := layoutUtils.GetRepositories()
	if err != nil {
		r.log.Error().Err(err).Msg("unable to get repositories list")

		return nil, err
	}

	if len(repoList) == 0 {
		r.log.Info().Msg("no repositories found")

		return imageList, nil
	}

	imageDir, imageTag := common.GetImageDirAndTag(image)

	imageManifest, _, err := layoutUtils.GetImageManifest(imageDir, imageTag)
	if err != nil {
		r.log.Info().Str("image", image).Msg("image not found")

		return imageList, err
	}

	imageLayers := imageManifest.Layers

	// This logic may not scale well in the future as we need to read all the
	// manifest files from the disk when the call is made, we should improve in a future PR
	for _, repo := range repoList {
		repoInfo, err := r.ExpandedRepoInfo(ctx, repo)
		if err != nil {
			r.log.Error().Err(err).Msg("unable to get image list")

			return nil, err
		}

		imageSummaries := repoInfo.Images

		var addImageToList bool
		// verify every image
		for _, imageSummary := range imageSummaries {
			if imageTag == *imageSummary.Tag && imageDir == repo {
				continue
			}

			addImageToList = true
			layers := imageSummary.Layers

			for _, l := range layers {
				foundLayer := false

				for _, k := range imageLayers {
					if *l.Digest == k.Digest.String() {
						foundLayer = true

						break
					}
				}

				if !foundLayer {
					addImageToList = false

					break
				}
			}

			if addImageToList {
				imageList = append(imageList, imageSummary)
			}
		}
	}

	return imageList, nil
}

// Image is the resolver for the Image field.
func (r *queryResolver) Image(ctx context.Context, image string) (*gql_generated.ImageSummary, error) {
	repo, tag := common.GetImageDirAndTag(image)

	return getImageSummary(ctx, repo, tag, r.repoDB, r.cveInfo, r.log)
}

// Referrers is the resolver for the Referrers field.
func (r *queryResolver) Referrers(ctx context.Context, repo string, digest string, typeArg string) ([]*gql_generated.Referrer, error) {
	store := r.storeController.GetImageStore(repo)

	referrers, err := getReferrers(store, repo, digest, typeArg, r.log)
	if err != nil {
		r.log.Error().Err(err).Msg("unable to get referrers from default store")

		return []*gql_generated.Referrer{}, err
	}

	return referrers, nil
}

// Query returns gql_generated.QueryResolver implementation.
func (r *Resolver) Query() gql_generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
