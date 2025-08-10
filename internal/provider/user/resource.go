// Package user provides the implementation of the User resource.
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tryretool/terraform-provider-retool/internal/provider/utils"
	"github.com/tryretool/terraform-provider-retool/internal/sdk/api"
)

// Ensure userResource implements the tfsdk.Resource interface.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// userResource schema structure.
type userResource struct {
	client *api.APIClient
}

// userResourceModel defines the data model for the User resource.
type userResourceModel struct {
	ID         types.String `tfsdk:"id"`
	LegacyID   types.String `tfsdk:"legacy_id"`
	Email      types.String `tfsdk:"email"`
	Active     types.Bool   `tfsdk:"active"`
	CreatedAt  types.String `tfsdk:"created_at"`
	LastActive types.String `tfsdk:"last_active"`
	FirstName  types.String `tfsdk:"first_name"`
	LastName   types.String `tfsdk:"last_name"`
	Metadata   types.Map    `tfsdk:"metadata"`
}

// Create new User resource.
func NewResource() resource.Resource {
	return &userResource{}
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*utils.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *utils.ProviderData, got: %T. Please report this issue to the provider developers.",
		)
	}
	r.client = providerData.Client
}

// Metadata associated with the User resource.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema returns the schema for the User resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `A user in Retool.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the user. Currently this is the same legacy_id but will change in the future.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"legacy_id": schema.StringAttribute{
				Computed:    true,
				Description: "The legacy ID of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "The email of the user.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the user is active or not. Defaults to true.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the user was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_active": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the user was last active.",
			},
			"first_name": schema.StringAttribute{
				Optional:    true,
				Description: "The first name of the user.",
			},
			"last_name": schema.StringAttribute{
				Optional:    true,
				Description: "The last name of the user.",
			},
			"metadata": schema.MapAttribute{
				Optional:    true,
				Description: "Metadata associated with the user.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Create creates the user resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	metadata := make(map[string]interface{}, len(plan.Metadata.Elements()))
	diag := plan.Metadata.ElementsAs(ctx, metadata, false)
	if diag.HasError() {
		resp.Diagnostics.AddError(
			"Error converting metadata",
			"Could not convert metadata elements to string",
		)
		return
	}

	// Generate API request body from plan.
	var user api.UsersPostRequest
	user.Email = plan.Email.ValueString()
	user.FirstName = plan.FirstName.ValueString()
	user.LastName = plan.LastName.ValueString()
	if plan.Active.IsNull() || !plan.Active.ValueBool() {
		myBool := true
		user.Active = &myBool
	} else {
		user.Active = plan.Active.ValueBoolPointer()
	}
	user.Metadata = metadata

	tflog.Info(ctx, "Creating a user", map[string]interface{}{"email": plan.Email.ValueString()})

	// Create new user.
	response, httpResponse, err := r.client.UsersAPI.UsersPost(ctx).UsersPostRequest(user).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		tflog.Error(ctx, "Error creating user", utils.AddHTTPStatusCode(map[string]interface{}{"error": err.Error()}, httpResponse))
		return
	}

	if response.Data.Id == "" {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: missing ID or LegacyID",
		)
		tflog.Error(ctx, "Error creating user", map[string]interface{}{"error": "missing ID or LegacyID"})
		return
	}

	// Map response body to schema and populate Computed attribute values.
	plan.ID = types.StringValue(response.Data.Id)
	plan.LegacyID = types.StringValue(utils.Float32PtrToIntString(&response.Data.LegacyId))

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error creating user", map[string]interface{}{"error": "Could not set state"})
		return
	}

	tflog.Info(ctx, "User created", map[string]interface{}{"id": response.Data.Id, "legacyId": response.Data.LegacyId})
}

// Read a User resource.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use the ID from the state to read the user.
	userID := state.ID.ValueString()
	user, httpResponse, err := r.client.UsersAPI.UsersUserIdGet(ctx, userID).Execute()
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			tflog.Info(ctx, "User not found", map[string]any{"userID": userID})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user with ID %s: %s", userID, err.Error()),
		)
		tflog.Error(ctx, "Error reading user", utils.AddHTTPStatusCode(map[string]any{"userID": userID, "error": err.Error()}, httpResponse))
		return
	}

	// Map the API response to the Terraform state.
	state.ID = types.StringValue(user.Data.Id)
	state.LegacyID = types.StringValue(utils.Float32PtrToIntString(&user.Data.LegacyId))
	state.Email = types.StringValue(user.Data.Email)
	state.Active = types.BoolValue(user.Data.Active)
	state.CreatedAt = types.StringValue(time.Time.String(user.Data.CreatedAt))
	if user.Data.LastActive.IsSet() {
		state.LastActive = types.StringValue(user.Data.LastActive.Get().Format(time.RFC3339))
	} else {
		state.LastActive = types.StringNull()
	}
	if user.Data.FirstName.IsSet() {
		state.FirstName = types.StringValue(*user.Data.FirstName.Get())
	} else {
		state.FirstName = types.StringNull()
	}
	if user.Data.LastName.IsSet() {
		state.LastName = types.StringValue(*user.Data.LastName.Get())
	} else {
		state.LastName = types.StringNull()
	}

	metadataVal, diag := types.MapValueFrom(ctx, types.StringType, user.Data.Metadata)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}
	state.Metadata = metadataVal

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a User resource. Currently not supported
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	return

}

// Delete a User resource.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := state.ID.ValueString()
	httpResponse, err := r.client.UsersAPI.UsersUserIdDelete(ctx, userID).Execute()
	if err != nil && !(httpResponse != nil && httpResponse.StatusCode == 404) { // It's ok to not find the user when deleting.
		resp.Diagnostics.AddError(
			"Error Deleting User",
			"Could not delete user with ID "+userID+": "+err.Error(),
		)
		tflog.Error(ctx, "Error Deleting User", utils.AddHTTPStatusCode(map[string]any{"error": err.Error(), "userID": userID}, httpResponse))
		return
	}
}

// ImportState allows importing of a User resource.
func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
