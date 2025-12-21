package preset

type ProcessStats struct {
	Total       int
	Success     int
	Failed      int
	WithErrors  int
	TotalErrors int
}

func GetStats(processed []ProcessedEntity) ProcessStats {
	var stats ProcessStats
	stats.Total = len(processed)

	for _, p := range processed {
		if p.FatalError != nil {
			stats.Failed++
		} else {
			stats.Success++
		}
		if len(p.Errors) > 0 {
			stats.WithErrors++
			stats.TotalErrors += len(p.Errors)
		}
	}
	return stats
}

func HasValidationErrors(processed []ProcessedEntity) bool {
	for _, p := range processed {
		if len(p.Errors) > 0 {
			return true
		}
	}
	return false
}

func CollectValidationErrors(processed []ProcessedEntity) map[string][]string {
	errs := make(map[string][]string)
	for _, p := range processed {
		if len(p.Errors) > 0 {
			errs[p.File.Path] = p.Errors
		}
	}
	return errs
}
