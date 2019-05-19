package cmd


const exitCodeUnexpected = 1
const exitCodeInvalidStatus = 2
const exitCodeInvalidResponse = 3
const exitCodeDryrun = 4
const exitCodeInvalidFlags = 5


func getOutputFormat(outputFormat string, defaultOutputFormat string) string {
	if outputFormat == "" {
		outputFormat = format
		if outputFormat == "" {
			outputFormat = defaultOutputFormat
		}
	}
	return outputFormat
}
