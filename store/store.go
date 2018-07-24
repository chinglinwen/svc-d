package store

type User struct {
    Name string
    Age  int
    Location string
    Password string
    Address string 
}

func (user *User) save(db *bolt.DB) error {
    // Store the user model in the user bucket using the username as the key.
    err := db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists(usersBucket)
        if err != nil {
            return err
        }    

        encoded, err := json.Marshal(user)
        if err != nil {
            return err
        }
        return b.Put([]byte(user.Name), encoded)
    })
    return err
}