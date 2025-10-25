package theme

type Options struct {
	Title      string
	Categories []Category
}

type Category struct {
	Name, URL string
}

func NewOptions() *Options {
	return &Options{
		Title: "Platform",
		Categories: []Category{
			{
				Name: "Category 1",
				URL:  "/category-1",
			},
		},
	}
}
