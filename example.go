package defaults

import "fmt"

func HasDefaults(values ...any) {
	first, isInt := DefaultValue(3).SafeCheck(values, 0, "First value must be an int")
	second, isString := DefaultValue("default").SafeCheck(values, 1, "Second value must be a string")

	if err := CheckDefaults(isInt, isString); err != nil {
		panic(err)
	}

	fmt.Println(first, second)
}
