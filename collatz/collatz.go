// Package collatz contains the functions and data pertaining to the collatz conjecture
package collatz

import (
	"fmt"
)

type Bud struct {
	Value    int
	Parent   *Twig
	Children []*Twig
}

type Twig struct {
	Parent *Bud
	Child  *Bud
	XAngle int
	YAngle int
}

type OrganicTree struct {
	Root *Bud
	Buds map[int]*Bud
}

func NextInt(current int) int {
	if current%2 == 0 {
		return current / 2
	} else {
		return current*3 + 1
	}
}

func BuildTree(maxNumber int) OrganicTree {
	if maxNumber < 1 {
		maxNumber = 1000
	}
	maxNumber++

	rootBud := Bud{
		Value:    1,
		Children: []*Twig{},
	}
	rec := OrganicTree{
		Root: &rootBud,
		Buds: map[int]*Bud{},
	}
	rec.Buds[1] = &rootBud

	for i := range maxNumber {
		if i == 0 || i == 1 {
			continue
		}

		bud := Bud{
			Value:    i,
			Children: []*Twig{},
		}
		rec.Buds[i] = &bud
		innerGrow := i
		for {
			nextInt := NextInt(innerGrow)
			parent, exists := rec.Buds[nextInt]
			if exists {
				xAngle, yAngle := getAngle(nextInt)
				newTwig := Twig{
					Child:  &bud,
					Parent: parent,
					XAngle: xAngle,
					YAngle: yAngle,
				}
				bud.Parent = &newTwig
				parent.Children = append(parent.Children, &newTwig)
				break
			} else {
				innerGrow = nextInt
				newBud := Bud{
					Value:    innerGrow,
					Children: []*Twig{},
				}
				rec.Buds[innerGrow] = &newBud
			}
		}

	}

	return rec
}

func PrintOrganicTree(t *OrganicTree) {
	if t.Root == nil {
		fmt.Println("(empty tree)")
		return
	}
	printBud(t.Root, "", true)
}

func printBud(b *Bud, prefix string, isLast bool) {
	connector := "├─"
	nextPrefix := prefix + "│  "
	if isLast {
		connector = "└─"
		nextPrefix = prefix + "   "
	}

	// Print this Bud
	if prefix == "" {
		fmt.Printf("Bud(%d)\n", b.Value)
	} else {
		fmt.Printf("%s%s Bud(%d)\n", prefix, connector, b.Value)
	}

	// Print its Twigs → Children
	for i, twig := range b.Children {
		lastTwig := i == len(b.Children)-1

		twigConnector := "├─"
		twigNextPrefix := nextPrefix + "│  "
		if lastTwig {
			twigConnector = "└─"
			twigNextPrefix = nextPrefix + "   "
		}

		fmt.Printf("%s%s Twig(x=%d, y=%d) → Bud(%d)\n",
			nextPrefix, twigConnector, twig.XAngle, twig.YAngle, twig.Child.Value)

		// Recursively print the child Bud
		printBud(twig.Child, twigNextPrefix, lastTwig)
	}
}

func getAngle(input int) (int, int) {
	if input%2 == 0 {
		return -8, 3
	} else {
		return 13, 3
	}
}
