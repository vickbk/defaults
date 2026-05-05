package defaults

import "fmt"

func HasDefaults(values ...any) {
	first, isInt := DefaultValue(3).SafeCheck(values, 0)
	second, isString := DefaultValue("default").SafeCheck(values, 1)

	if err := CheckDefaults(isInt, "First value must be an int", isString, "second value must be string"); err != nil {
		panic(err)
	}

	fmt.Println(first, second)
}
