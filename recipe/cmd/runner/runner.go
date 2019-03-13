package main

import (
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/recipe/cmd/commons"
	"fmt"
	"os"
)

func main() {
	commander := &recipe.IOCommander{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	packsConf := commons.PacksConfig()
	executor := &recipe.PacksExecutor{
		Conf:      packsConf,
		Commander: commander,
	}

	recipeConf := commons.RecipeConfig()
	err := executor.ExecuteRecipe(recipeConf)
	if err != nil {
		commons.RespondWithFailure(err)
		os.Exit(1)
	}

	fmt.Println("Recipe Execution completed")
}
