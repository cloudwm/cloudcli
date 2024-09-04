package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func getInteractiveFlagLongName(flag SchemaCommandFlag) string {
	if flag.LongName != "" {
		return flag.LongName
	} else {
		return flag.Name
	}
}

func getInteractiveFlagValue_Datacenter(reader *bufio.Reader) string {
	respString := getListOfListsRespString("cloudcli-server-options.json", true, "/service/server")
	rootitems := jsonUnmarshalItemsList(respString)
	var datacenter_keys []string
	var datacenter_descriptions []string
	for rootlistkey, rootitem := range rootitems {
		if rootlistkey == "datacenters" {
			for itemkey, itemvalue := range rootitem.(map[string]interface{}) {
				datacenter_keys = append(datacenter_keys, parseItemString(itemkey))
				datacenter_descriptions = append(datacenter_descriptions, parseItemString(itemvalue))
			}
			break
		}
	}
	sort.Strings(datacenter_descriptions)
	fmt.Printf("%s\n", strings.Join(datacenter_descriptions, "\n"))
	selected_datacenter := ""
	for selected_datacenter == "" {
		fmt.Printf("Choose a datacenter from the list, type the datacenter letters code: ")
		selected_datacenter = readInput(reader)
		if selected_datacenter == "" {
			fmt.Printf("Datacenter is required. ")
		} else {
			for _, datacenter_key := range datacenter_keys {
				if datacenter_key == selected_datacenter {
					return selected_datacenter
				}
			}
			fmt.Printf("Invalid value. ")
			selected_datacenter = ""
		}
	}
	return ""
}

func getCpuLetter(cpu string) string {
	for _, letter := range []string{"A", "T", "D", "B"} {
		if strings.HasSuffix(cpu, letter) {
			return letter
		}
	}
	return ""
}

func getCpuLetterCategory(cpuLetter string) string {
	if cpuLetter == "A" {
		return "Availability"
	} else if cpuLetter == "D" {
		return "Dedicated"
	} else if cpuLetter == "T" {
		return "Burstable"
	} else if cpuLetter == "B" {
		return "General"
	} else {
		return ""
	}
}

func getInteractiveFlagValue_Ram(reader *bufio.Reader, defaultRam string, cpu string) string {
	respString := getListOfListsRespString("cloudcli-server-options.json", true, "/service/server")
	rootitems := jsonUnmarshalItemsList(respString)
	cpuLetter := getCpuLetter(cpu)
	var rams []float64
	for rootlistkey, rootitem := range rootitems {
		if rootlistkey == "ram" {
			for itemkey, itemvalues := range rootitem.(map[string]interface{}) {
				if itemkey == cpuLetter {
					for _, ram := range itemvalues.([]interface{}) {
						rams = append(rams, ram.(float64))
					}
				}
			}
		}
	}
	sort.Float64s(rams)
	for _, ram := range rams {
		fmt.Printf("%d\n", int(ram))
	}
	ok := false
	selectedRam := ""
	for !ok {
		fmt.Printf("Enter a RAM value from the list (default=%s): ", defaultRam)
		selectedRam = readInput(reader)
		if selectedRam == "" {
			selectedRam = defaultRam
		}
		for _, ram := range rams {
			if ram == cast.ToFloat64(selectedRam) {
				ok = true
			}
		}
		if !ok {
			fmt.Printf("Invalid RAM value\n")
		}
	}
	return selectedRam
}

func getInteractiveFlagValue_Cpu(reader *bufio.Reader, defaultCpu string) string {
	respString := getListOfListsRespString("cloudcli-server-options.json", true, "/service/server")
	rootitems := jsonUnmarshalItemsList(respString)
	cpus := make(map[string][]string)
	for rootlistkey, rootitem := range rootitems {
		if rootlistkey == "cpu" {
			for _, itemvalue := range rootitem.([]interface{}) {
				letter := getCpuLetter(itemvalue.(string))
				if letter == "" {
					fmt.Printf("Invalid CPU: %s\n", itemvalue.(string))
					os.Exit(exitCodeUnexpected)
				}
				cpus[letter] = append(cpus[letter], itemvalue.(string))
			}
			break
		}
	}
	var letters []string
	for letter, cpus := range cpus {
		sort.Strings(cpus)
		letters = append(letters, letter)
	}
	sort.Strings(letters)
	for _, letter := range letters {
		category := getCpuLetterCategory(letter)
		if category == "" {
			fmt.Printf("Invalid CPU type: %s\n", letter)
			os.Exit(exitCodeUnexpected)
		}
		fmt.Printf("%s: %s\n", letter, category)
	}
	var defaultLetter string
	var defaultCpuCores string
	if len(defaultCpu) == 2 {
		defaultLetter = strings.Split(defaultCpu, "")[1]
		defaultCpuCores = strings.Split(defaultCpu, "")[0]
	} else if len(defaultCpu) == 3 {
		defaultLetter = strings.Split(defaultCpu, "")[2]
		defaultCpuCores = strings.Split(defaultCpu, "")[0] + strings.Split(defaultCpu, "")[1]
	} else {
		fmt.Printf("Invalid default CPU: %s\n", defaultCpu)
		os.Exit(exitCodeUnexpected)
	}
	ok := false
	selectedLetter := ""
	for !ok {
		fmt.Printf("Enter a CPU type from the list (default=%s): ", defaultLetter)
		selectedLetter = readInput(reader)
		if selectedLetter == "" {
			selectedLetter = defaultLetter
		}
		for _, letter := range letters {
			if letter == selectedLetter {
				ok = true
			}
		}
		if !ok {
			fmt.Printf("Invalid CPU type\n")
		}
	}
	fmt.Printf("\n")
	for _, cpu := range cpus[selectedLetter] {
		if len(cpu) < 3 {
			fmt.Printf("%s\n", strings.TrimSuffix(cpu, selectedLetter))
		}
	}
	for _, cpu := range cpus[selectedLetter] {
		if len(cpu) >= 3 {
			fmt.Printf("%s\n", strings.TrimSuffix(cpu, selectedLetter))
		}
	}
	ok = false
	selectedCpuCores := ""
	for !ok {
		fmt.Printf("Enter the number of CPU cores from the list (default=%s): ", defaultCpuCores)
		selectedCpuCores = readInput(reader)
		if selectedCpuCores == "" {
			selectedCpuCores = defaultCpuCores
		}
		for _, cpu := range cpus[selectedLetter] {
			if cpu == selectedCpuCores+selectedLetter {
				ok = true
			}
		}
		if !ok {
			fmt.Printf("Invalid CPU cores\n")
		}
	}
	return selectedCpuCores + selectedLetter
}

func getInteractiveFlagValue_Networks(reader *bufio.Reader, datacenter string) string {
	respString := getListOfListsRespString("cloudcli-server-options.json", true, "/service/server")
	rootitems := jsonUnmarshalItemsList(respString)
	var netnames []string
	for rootlistkey, rootitem := range rootitems {
		if rootlistkey == "networks" {
			for itemkey, itemvalue := range rootitem.(map[string]interface{}) {
				if itemkey == datacenter {
					for _, network := range itemvalue.([]interface{}) {
						netname := ""
						for netkey, netval := range network.(map[string]interface{}) {
							if netkey == "name" {
								netname = netval.(string)
							}
						}
						netnames = append(netnames, netname)
					}
					break
				}
			}
		}
	}
	fmt.Printf("Available network names: \n")
	for _, netname := range netnames {
		fmt.Printf("%s\n", netname)
	}
	fmt.Printf("You can assign up to 4 network interfaces from above network names\n")
	var selectedNetworks []string
	for netId, netNum := range []int{1, 2, 3, 4} {
		var selectedNetName string
		var selectedNetIp string
		ok := false
		for !ok {
			if netId == 0 {
				fmt.Printf("Network 1 name (default=wan): ")
				selectedNetName = readInput(reader)
				if selectedNetName == "" {
					selectedNetName = "wan"
				}
			} else {
				fmt.Printf("Network %d name (leave empty to stop adding networks): ", netNum)
				selectedNetName = readInput(reader)
			}
			if selectedNetName == "" {
				break
			} else {
				for _, netname := range netnames {
					if cast.ToString(selectedNetName) == netname {
						ok = true
					}
				}
				if !ok {
					fmt.Printf("Invalid network name, choose from the list\n")
				}
			}
		}
		if selectedNetName == "" {
			break
		}
		if selectedNetName == "wan" {
			selectedNetIp = "auto"
		} else {
			fmt.Printf("Enter the network interface IP (leave empty for auto): ")
			selectedNetIp = readInput(reader)
			if selectedNetIp == "" {
				selectedNetIp = "auto"
			}
		}
		selectedNetworks = append(selectedNetworks, fmt.Sprintf("id=%d,name=%s,ip=%s", netId, selectedNetName, selectedNetIp))
	}
	return strings.Join(selectedNetworks, " ")
}

func getInteractiveFlagValue_Disk(reader *bufio.Reader) string {
	respString := getListOfListsRespString("cloudcli-server-options.json", true, "/service/server")
	rootitems := jsonUnmarshalItemsList(respString)
	var diskSizes []float64
	for rootlistkey, rootitem := range rootitems {
		if rootlistkey == "disk" {
			for _, itemvalue := range rootitem.([]interface{}) {
				diskSizes = append(diskSizes, itemvalue.(float64))
			}
			break
		}
	}
	sort.Float64s(diskSizes)
	fmt.Printf("Available disk sizes:\n")
	for _, diskSize := range diskSizes {
		fmt.Printf("%d\n", int(diskSize))
	}
	fmt.Printf("You can create up to 4 disks with sizes from above list (in GB)\n")
	var selectedSizes []string
	for diskId, diskNum := range []int{1, 2, 3, 4} {
		var size string
		ok := false
		for !ok {
			if diskId == 0 {
				fmt.Printf("Disk 1 GB size (default=20): ")
				size = readInput(reader)
				if size == "" {
					size = "20"
				}
			} else {
				fmt.Printf("Disk %d GB size (leave empty to stop adding disks): ", diskNum)
				size = readInput(reader)
			}
			if size == "" {
				break
			} else {
				for _, diskSize := range diskSizes {
					if cast.ToString(diskSize) == size {
						ok = true
					}
				}
				if !ok {
					fmt.Printf("Invalid size, choose from the list\n")
				}
			}
		}
		if size == "" {
			break
		}
		selectedSizes = append(selectedSizes, fmt.Sprintf("id=%d,size=%s", diskId, size))
	}
	return strings.Join(selectedSizes, " ")
}

func getInteractiveFlagValue_Traffic(reader *bufio.Reader, datacenter string) string {
	respString := getListOfListsRespString("", false, "/svc?path=serverCreate/datacenterConfiguration/"+datacenter)
	rootitems := jsonUnmarshalItemsList(respString)
	var trafficPackageOptions []string
	defaultTrafficPackage := ""
	for rootlistkey, rootitem := range rootitems {
		if rootlistkey == "trafficPackage" {
			defaultTrafficPackage = rootitem.(string)
		} else if rootlistkey == "trafficPackageConf" {
			for _, opt := range rootitem.([]interface{}) {
				trafficPackageOptions = append(trafficPackageOptions, opt.(string))
			}
		}
	}
	selectedTrafficOpt := ""
	for true {
		fmt.Printf("Enter a traffic option (%s) (default=%s): ", strings.Join(trafficPackageOptions, "|"), defaultTrafficPackage)
		selectedTrafficOpt = readInput(reader)
		if selectedTrafficOpt == "" {
			selectedTrafficOpt = defaultTrafficPackage
		}
		ok := false
		for _, opt := range trafficPackageOptions {
			if selectedTrafficOpt == opt {
				ok = true
			}
		}
		if ok {
			break
		} else {
			fmt.Printf("Invalid option. ")
		}
	}
	return selectedTrafficOpt
}

func getInteractiveFlagValue_Image(reader *bufio.Reader, datacenter string) string {
	respString := getListOfListsRespString("cloudcli-server-options.json", true, "/service/server")
	rootitems := jsonUnmarshalItemsList(respString)
	root_images := make(map[string][]string)
	for rootlistkey, rootitem := range rootitems {
		if rootlistkey == "diskImages" {
			for itemkey, itemvalue := range rootitem.(map[string]interface{}) {
				if itemkey == datacenter {
					for _, raw_image := range itemvalue.([]interface{}) {
						image := raw_image.(map[string]interface{})
						splitimage := strings.Split(image["description"].(string), "_")
						root_images[splitimage[0]] = append(root_images[splitimage[0]], strings.Join(splitimage, "_"))
					}
				}
			}
			break
		}
	}
	var categories []string
	for category, images := range root_images {
		sort.Strings(images)
		if category != "service" && category != "apps" {
			categories = append(categories, category)
		}
	}
	sort.Strings(categories)
	for _, category := range categories {
		fmt.Printf("%s\n", category)
	}
	selectedCategory := ""
	for selectedCategory == "" {
		fmt.Printf("Enter an image type: ")
		selectedCategory = readInput(reader)
		if selectedCategory == "" {
			continue
		}
		ok := false
		for _, category := range categories {
			if selectedCategory == category {
				ok = true
			}
		}
		if ok {
			break
		}
		fmt.Printf("Invalid image type\n")
		selectedCategory = ""
	}
	fmt.Printf("\n")
	for _, image := range root_images[selectedCategory] {
		fmt.Printf("%s\n", image)
	}
	selectedImage := ""
	for selectedImage == "" {
		fmt.Printf("Enter an image: ")
		selectedImage = readInput(reader)
		if selectedImage == "" {
			continue
		}
		ok := false
		for _, image := range root_images[selectedCategory] {
			if selectedImage == image {
				ok = true
			}
		}
		if ok {
			break
		}
		fmt.Printf("Invalid image\n")
		selectedImage = ""
	}
	return selectedImage
}

func readInput(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	text = strings.TrimSpace(text)
	return text
}

func getInteractiveFlagValue(flag SchemaCommandFlag, reader *bufio.Reader, datacenter string, cpu string) string {
	fmt.Printf("\n")
	if flag.SelectfromServeroption == "datacenter" {
		return getInteractiveFlagValue_Datacenter(reader)
	} else if flag.SelectfromServeroption == "image" {
		return getInteractiveFlagValue_Image(reader, datacenter)
	} else if flag.SelectfromServeroption == "cpu" {
		return getInteractiveFlagValue_Cpu(reader, flag.Default)
	} else if flag.SelectfromServeroption == "ram" {
		return getInteractiveFlagValue_Ram(reader, flag.Default, cpu)
	} else if flag.SelectfromServeroption == "disk" {
		return getInteractiveFlagValue_Disk(reader)
	} else if flag.SelectfromServeroption == "network" {
		return getInteractiveFlagValue_Networks(reader, datacenter)
	} else if flag.SelectfromServeroption == "traffic" {
		return getInteractiveFlagValue_Traffic(reader, datacenter)
	} else {
		if flag.Name == "ssh-key" {
			fmt.Printf("Absolute path to public key file, adds to server authorized keys after create is done: ")
		} else {
			fmt.Printf("%s: ", flag.Usage)
		}
		text := readInput(reader)
		if flag.Required && text == "" {
			fmt.Printf("%s is required, please enter a value\n", getInteractiveFlagLongName(flag))
			return getInteractiveFlagValue(flag, reader, datacenter, cpu)
		}
		if flag.ValidateRegex != "" {
			if matched, err := regexp.MatchString("^"+flag.ValidateRegex+"$", text); err != nil || !matched {
				fmt.Printf("%s must match regular expression: '%s'\n", getInteractiveFlagLongName(flag), flag.ValidateRegex)
				return getInteractiveFlagValue(flag, reader, datacenter, cpu)
			}
		}
		if flag.ValidatePassword {
			for _, pattern := range []string{"a-z", "A-Z", "0-9"} {
				if matched, err := regexp.MatchString(fmt.Sprintf("[%s]", pattern), text); err != nil || !matched {
					fmt.Printf("%s must contain at least one character in range [%s]", getInteractiveFlagLongName(flag), pattern)
					return getInteractiveFlagValue(flag, reader, datacenter, cpu)
				}
			}
		}
		if flag.ValidateBoolean {
			if text != "yes" && text != "no" && text != "" {
				fmt.Printf("Please enter 'yes' or 'no'. ")
				return getInteractiveFlagValue(flag, reader, datacenter, cpu)
			}
		}
		if flag.ValidateIntegerMin != 0 && text != "" {
			if num, err := strconv.Atoi(text); err != nil || num < flag.ValidateIntegerMin {
				fmt.Printf("%s must be at least %d. ", getInteractiveFlagLongName(flag), flag.ValidateIntegerMin)
				return getInteractiveFlagValue(flag, reader, datacenter, cpu)
			}
		}
		if flag.ValidateIntegerMax != 0 && text != "" {
			if num, err := strconv.Atoi(text); err != nil || num > flag.ValidateIntegerMax {
				fmt.Printf("%s must be less then or equal to %d. ", getInteractiveFlagLongName(flag), flag.ValidateIntegerMax)
				return getInteractiveFlagValue(flag, reader, datacenter, cpu)
			}
		}
		if len(flag.ValidateValues) > 0 && text != "" {
			ok := false
			for _, value := range flag.ValidateValues {
				if text == value {
					ok = true
				}
			}
			if !ok {
				fmt.Printf("Invalid value for %s. ", getInteractiveFlagLongName(flag))
				return getInteractiveFlagValue(flag, reader, datacenter, cpu)
			}
		}
		return text
	}
}

func commandRunPostInteractive(cmd *cobra.Command, command SchemaCommand) {
	fmt.Printf("Fetching server options... ")
	refreshListOfListsCache("cloudcli-server-options.json")
	fmt.Printf("OK\n")
	reader := bufio.NewReader(os.Stdin)
	datacenter := ""
	cpu := ""
	for _, flag := range command.Flags {
		if flag.Name == "interactive" {
			continue
		}
		if flag.Name == "wait" {
			_ = cmd.Flags().Set("wait", "1")
			continue
		}
		value := getInteractiveFlagValue(flag, reader, datacenter, cpu)
		if flag.Name == "datacenter" {
			datacenter = value
		} else if flag.Name == "cpu" {
			cpu = value
		}
		_ = cmd.Flags().Set(flag.Name, value)
	}
	fmt.Printf("\nCreating server..\n")
}
