package defaults

import "fmt"

func HasDefaults(values ...any) {
	first, isInt := Value(3).SafeCheck(values, 0, "First value must be an int")
	second, isString := Value("default").SafeCheck(values, 1, "Second value must be a string")

	if err := AggregateErrors(isInt, isString); err != nil {
		panic(err)
	}

	fmt.Println(first, second, isInt.Ok, isString.Ok, isInt.UsedDefault, isString.UsedDefault)
}
