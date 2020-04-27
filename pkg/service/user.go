package service

// BORRAR

func getConnections(eventSubdomain, audienceType string, getter connectionGetter) ([]string, error) {
	var connections []string

	if err := getter.GetUserConnections(eventSubdomain, audienceType, &connections); err != nil {
		return nil, err
	}

	return nil, nil
}
