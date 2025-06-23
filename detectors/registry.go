package detectors

// GetAllDetectors returns all available detectors
func GetAllDetectors() []Detector {
	return []Detector{
		&GitDetector{},
		&GitLabCIDetector{},
		&GitHubActionsDetector{},
		&CircleCIDetector{},
		&TravisCIDetector{},
		&AzurePipelinesDetector{},
		&JenkinsDetector{},
		&BitbucketPipelinesDetector{},
		&GemfileDetector{},
		&RequirementsDetector{},
		&PackageJSONDetector{},
		&ComposerDetector{},
		&GoModDetector{},
		&CargoDetector{},
		&DotNetDetector{},
		&JavaDetector{},
		&DartDetector{},
		&IOSDetector{},
		&AndroidDetector{},
		&ChromeExtensionDetector{},
		&VSCodeExtensionDetector{},
		&ShopifyDetector{},
		&NxDetector{},
		&WordPressDetector{},
		&ContainerRegistryDetector{},
		&I18nDetector{},
		&AIServicesDetector{},
		&SearchServicesDetector{},
		&MapsServicesDetector{},
		&PushNotificationsDetector{},
		&VercelDetector{},
		&NetlifyDetector{},
		&HerokuDetector{},
		&FirebaseDetector{},
		&RailwayDetector{},
		&RenderDetector{},
		&FlyIODetector{},
	}
}

// GetDetectorNames returns list of all detector names
func GetDetectorNames() []string {
	detectors := GetAllDetectors()
	names := make([]string, len(detectors))
	for i, detector := range detectors {
		names[i] = detector.Name()
	}
	return names
}

// FindDetectorByName finds detector by name
func FindDetectorByName(name string) Detector {
	for _, detector := range GetAllDetectors() {
		if detector.Name() == name {
			return detector
		}
	}
	return nil
}