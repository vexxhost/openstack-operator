package baseutils

// CompareStrSlice compares two string slices and return the different elements
// Return values are 2 arrays; aOnlySlice, and bOnlySlice
func CompareStrSlice(aS []string, bS []string) ([]string, []string) {
	aOnlyS := []string{}
	for _, a := range aS {
		i, isExist := Find(bS, a)
		if !isExist {
			aOnlyS = append(aOnlyS, a)
		} else {
			RemoveElement(&bS, i)
		}
	}
	return aOnlyS, bS
}

// Find is a helper function to find the string in a slice of strings.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// RemoveElement is a helper function to remove the ith string from a slice of strings.
func RemoveElement(a *[]string, i int) {
	(*a)[i] = (*a)[len(*a)-1] // Copy last element to index i.
	(*a)[len(*a)-1] = ""      // Erase last element (write zero value).
	(*a) = (*a)[:len(*a)-1]   // Truncate the length
}

// ContainsString is a helper function to check string in a slice of strings
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// RemoveString is a helper function to remove string from a slice of strings.
func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
