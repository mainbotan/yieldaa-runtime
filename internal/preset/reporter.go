package preset

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintResults(pkg *Package, processed []ProcessedEntity, fatalErrs []error) {
	// Header
	fmt.Printf("%s v%s\n", pkg.Name, pkg.Version)
	fmt.Printf("files:%d size:%.1fKB hash:%08x\n\n",
		pkg.EntitiesCount, float64(pkg.EntitiesTotalSize)/1024, pkg.EntitiesStructureHash)

	// Processing stats
	stats := GetStats(processed)
	fmt.Printf("processed:%d failed:%d errors:%d\n\n",
		stats.Success, stats.Failed, stats.TotalErrors)

	// Files table
	if len(processed) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "STATUS\tPATH\tSIZE\tCONTENT_HASH")
		fmt.Fprintln(w, "------\t----\t----\t-----------")

		for _, p := range processed {
			status := "✓"
			if p.FatalError != nil {
				status = "✗"
			} else if len(p.Errors) > 0 {
				status = "!"
			}

			path := ShortPath(p.File.Path, 40)
			size := fmt.Sprintf("%.1fK", float64(p.File.Size)/1024)
			hash := fmt.Sprintf("%08x", p.ContentHash)

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", status, path, size, hash)
		}
		w.Flush()
	}

	// Fatal errors
	if len(fatalErrs) > 0 {
		fmt.Printf("\nFATAL ERRORS:\n")
		for _, err := range fatalErrs {
			fmt.Printf("  %v\n", err)
		}
	}

	// Validation errors
	validationErrs := CollectValidationErrors(processed)
	if len(validationErrs) > 0 {
		fmt.Printf("\nVALIDATION ERRORS:\n")
		for path, errs := range validationErrs {
			fmt.Printf("  %s\n", ShortPath(path, 50))
			for _, err := range errs {
				fmt.Printf("    • %s\n", err)
			}
		}
	}
}
