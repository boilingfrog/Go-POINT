package main

func main() {
	cat := Cat{}

	bigCat := BigCat{&cat}
	bigCat.Eat()
}
