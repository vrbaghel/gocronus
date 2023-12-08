package env

func LoadENV() (*ENV, error) {
	env, err := readENV()
	if err != nil {
		return nil, err
	}

	if err = validateENVs(*env); err != nil {
		return nil, err
	}

	return env, nil
}
