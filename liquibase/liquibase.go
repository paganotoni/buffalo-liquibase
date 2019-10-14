package liquibase

const Version = "v0.0.1"

// func BuildRunArgsFor(environment string) ([]string, error) {
// 	env := pop.Connections[environment]
// 	if env == nil {
// 		return []string{}, fmt.Errorf("could not find %v environment in your database.yml", environment)
// 	}

// 	originalURL := env.URL()

// 	r := regexp.MustCompile(`postgres:\/\/(?P<username>.*):(?P<password>.*)@(?P<host>.*):(?P<port>.*)\/(?P<database>.*)\?(?P<extras>.*)`)
// 	match := r.FindStringSubmatch(originalURL)
// 	if match == nil {
// 		return []string{}, fmt.Errorf("could not convert %v url into liquibase", environment)
// 	}

// 	URL := fmt.Sprintf("jdbc:postgresql://%v:%v/%v?%v", match[3], match[4], match[5], match[6])
// 	runArgs := []string{
// 		"--driver=org.postgresql.Driver",
// 		"--url=" + URL,
// 		"--logLevel=info",
// 		"--username=" + match[1],
// 		"--password=" + match[2],
// 	}
// 	return runArgs, nil
// }
