package main


type Community struct {
    id int `json:"id"`
    name string `json:"name"`
    description string `json:"description"`
}

func (c *Community) SelectPosts(n int) []int {
    posts := []int{1, 2} // placeholder TODO remove me
    return posts
}

func (c *Community) Id() int {
    return c.id
}

func (c *Community) Name() string {
    return c.name
}

func (c *Community) Description() string {
    return c.description
}
