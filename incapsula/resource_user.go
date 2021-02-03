package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"time"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"email": {
				Description: "Email address. For example: joe@example.com. example: userEmail@imperva.com",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"account_id": {
				Description: "Unique ID of the required account . example: 123456",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"first_name": {
				Description: "The first name of the user that was acted on. example: John",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"last_name": {
				Description: "The last name of the user that was acted on. example: Snow",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"role_names": {
				Description: "List of role names to add to the user. Use roleIds or roleNames to add roles to the user, but not both.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	email := d.Get("email").(string)
	accountId := d.Get("account_id").(int)

	log.Printf("[INFO] Creating Incapsula user for email: %s\n", email)

	UserAddResponse, err := client.AddUser(
		accountId,
		email,
		d.Get("role_names").([]interface{}),
		d.Get("first_name").(string),
		d.Get("last_name").(string),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create user for email: %s, %s\n", email, err)
		return err
	}

	// Set the User ID
	d.SetId(strconv.Itoa(UserAddResponse.UserID))
	log.Printf("[INFO] Created Incapsula user for email: %s\n", email)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	// Set the rest of the state from the resource read
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	userID, _ := strconv.Atoi(d.Id())
	accountID := d.Get("account_id").(int)
	email := d.Get("email").(string)
	log.Printf("[INFO] Reading Incapsula user : %d\n", userID)

	UserStatusResponse, err := client.UserStatus(accountID, email)

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula user: %s, %s\n", email, err)
		return err
	}

	log.Printf("[INFO]listRoles : %v\n", UserStatusResponse.Roles)

	time.Sleep(5 * time.Second)

	listRoles := make([]interface{}, len(UserStatusResponse.Roles))
	for i, v := range UserStatusResponse.Roles {
		log.Printf("[INFO]listRoles : %v\n", UserStatusResponse.Roles)
		time.Sleep(3 * time.Second)
		listRoles[i] = v.RoleName
	}

	d.Set("email", UserStatusResponse.Email)
	d.Set("account_id", UserStatusResponse.AccountID)
	d.Set("first_name", UserStatusResponse.FirstName)
	d.Set("last_name", UserStatusResponse.LastName)
	d.Set("role_names", listRoles)

	log.Printf("[INFO] Finished reading Incapsula user: %s\n", email)

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	email := d.Get("email").(string)

	log.Printf("[INFO] Deleting Incapsula user: %s\n", email)

	err := client.DeleteUser(accountID, email)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula user: %s %s\n", email, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] Deleted Incapsula user: %s\n", email)

	return nil
}
