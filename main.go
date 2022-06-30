package Go_POINT

func UpdateUser(userId int) {

}

func DelUser(userId int) {

}

func Demo(userId, usertype int) error {
	if usertype == 1 {

	} else {

	}

}

type User struct {
	Name      string
	Age       int
	AvatarUrl string
}

type User1 struct {
	UserName      string
	UserAge       int
	UserAvatarUrl string
}

func sum(sil []*User, age int) int {
	count := 0
	if len(sil) == 0 || age == 0 {
		return count
	} else {
		for _, item := range sil {
			if item.Age > age {
				count++
				// .....
			}
		}
	}
	return count
}
