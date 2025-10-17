package plantstorage

func checkJSONBKeys(jsonB map[string]interface{}, mustBeKeys map[string]error) error {
	for key := range jsonB {
		if _, ok := mustBeKeys[key]; !ok {
			if err := mustBeKeys[key]; err != nil {
				return err
			}
		}
	}
	return nil
}
