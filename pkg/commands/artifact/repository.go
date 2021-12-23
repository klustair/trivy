package artifact

import (
	"context"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	"github.com/aquasecurity/fanal/cache"
	pkgReport "github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
)

func repositoryScanner(ctx context.Context, dir string, ac cache.ArtifactCache, lac cache.LocalArtifactCache,
	_ time.Duration, artifactOpt artifact.Option, scannerOpt config.ScannerOption) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeRepositoryScanner(ctx, dir, ac, lac, artifactOpt, scannerOpt)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a filesystem scanner: %w", err)
	}
	return s, cleanup, nil
}

// RepositoryRun runs scan on repository
func RepositoryRun(ctx *cli.Context) error {
	opt, err := initOption(ctx)
	if err != nil {
		return xerrors.Errorf("option error: %w", err)
	}

	// Do not scan OS packages
	opt.VulnType = []string{types.VulnTypeLibrary}

	// Disable the OS analyzers and individual package analyzers
	opt.DisabledAnalyzers = append(analyzer.TypeIndividualPkgs, analyzer.TypeOSes...)

	return Run(ctx.Context, opt, repositoryScanner, initFSCache)
}

func RepositoryRunLib(ctx context.Context, opt Option) (pkgReport.Report, error) {

	// Do not scan OS packages
	opt.VulnType = []string{types.VulnTypeLibrary}

	// Disable the OS analyzers and individual package analyzers
	opt.DisabledAnalyzers = append(analyzer.TypeIndividualPkgs, analyzer.TypeOSes...)

	return run(ctx, opt, repositoryScanner, initFSCache)
}
