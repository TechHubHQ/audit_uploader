package main

func getOptions(question []any) []string {
	if len(question) < 5 {
		return nil
	}

	questionType, ok := question[3].(float64)
	if !ok || int(questionType) != 2 {
		return nil
	}

	entryData, ok := question[4].([]any)
	if !ok || len(entryData) == 0 {
		return nil
	}

	firstEntry, ok := entryData[0].([]any)
	if !ok || len(firstEntry) < 2 {
		return nil
	}

	optionData, ok := firstEntry[1].([]any)
	if !ok {
		return nil
	}

	var options []string

	for _, item := range optionData {
		option, ok := item.([]any)
		if !ok || len(option) == 0 {
			continue
		}

		text, ok := option[0].(string)
		if !ok {
			continue
		}

		options = append(options, text)
	}

	return options
}
