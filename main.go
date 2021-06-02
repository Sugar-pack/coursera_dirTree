package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func createmargin(level int, levoflast int) string {
	startchar := "│"
	tab := "\t"
	margin := ""
	if levoflast == 0 {
		for i := 0; i < level; i++ {
			margin += startchar
			margin += tab
		}
	} else {
		for i := 0; i < level-levoflast; i++ {
			margin += startchar
			margin += tab
		}
		for i := 0; i < levoflast; i++ {
			margin += tab
		}

	}
	return margin
}

func dirTreewithlevels(out io.Writer, path string, files bool, level int, levoflast int) error {
	sep := "├───"

	loflt := levoflast
	fls, err := os.ReadDir(path)

	margin := createmargin(level, loflt)

	if err != nil {
		return err
	}

	fls2 := make([]os.DirEntry, 0, len(fls))
	if !files {
		for _, unit := range fls {
			if unit.IsDir() {
				fls2 = append(fls2, unit)

			}
		}
		fls = fls2
	}

	sort.SliceStable(fls, func(i, j int) bool { return fls[i].Name() < fls[j].Name() })

	for i, unit := range fls {
		if i == len(fls)-1 {
			sep = "└───"
			loflt += 1
		}
		if unit.IsDir() {
			fmt.Fprintf(out, "%v%v%v\n", margin, sep, unit.Name())
			newpath := path + string(os.PathSeparator) + unit.Name()
			newerr := dirTreewithlevels(out, newpath, files, level+1, loflt)
			if newerr != nil {
				return newerr
			}
		} else {

			info, newerr := unit.Info()
			sizeoffile := info.Size()
			if newerr != nil {
				return newerr
			}
			if sizeoffile == 0 {
				fmt.Fprintf(out, "%v%v%v%v\n", margin, sep, unit.Name(), " (empty)")
			} else {
				fmt.Fprintf(out, "%v%v%v%v%v%v\n", margin, sep, unit.Name(), " (", sizeoffile, "b)")
			}

		}

	}
	return nil
}

func dirTree(out io.Writer, path string, files bool) error {
	sep := "├───"
	loflt := 0
	fls, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if !files {
		for i, unit := range fls {
			if !unit.IsDir() {
				fls[i] = fls[len(fls)-1]
				fls = fls[:len(fls)-1]
			}
		}
	}

	sort.SliceStable(fls, func(i, j int) bool { return fls[i].Name() < fls[j].Name() })

	for i, unit := range fls {
		if i == len(fls)-1 {
			sep = "└───"
			loflt = 1
		}
		if unit.IsDir() {
			fmt.Fprintf(out, "%v%v\n", sep, unit.Name())
			newpath := path + string(os.PathSeparator) + unit.Name()
			newerr := dirTreewithlevels(out, newpath, files, 1, loflt)
			if newerr != nil {
				return newerr
			}
		} else {

			info, newerr := unit.Info()
			sizeoffile := info.Size()
			if newerr != nil {
				return newerr
			}
			if sizeoffile == 0 {
				fmt.Fprintf(out, "%v%v%v\n", sep, unit.Name(), " (empty)")
			} else {
				fmt.Fprintf(out, "%v%v%v%v%v\n", "├───", unit.Name(), " (", sizeoffile, "b)")
			}

		}
	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
