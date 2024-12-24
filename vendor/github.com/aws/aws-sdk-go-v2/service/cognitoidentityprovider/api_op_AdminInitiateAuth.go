// Code generated by smithy-go-codegen DO NOT EDIT.

package cognitoidentityprovider

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Starts sign-in for applications with a server-side component, for example a
// traditional web application. This operation specifies the authentication flow
// that you'd like to begin. The authentication flow that you specify must be
// supported in your app client configuration. For more information about
// authentication flows, see [Authentication flows].
//
// This action might generate an SMS text message. Starting June 1, 2021, US
// telecom carriers require you to register an origination phone number before you
// can send SMS messages to US phone numbers. If you use SMS text messages in
// Amazon Cognito, you must register a phone number with [Amazon Pinpoint]. Amazon Cognito uses the
// registered number automatically. Otherwise, Amazon Cognito users who must
// receive SMS messages might not be able to sign up, activate their accounts, or
// sign in.
//
// If you have never used SMS text messages with Amazon Cognito or any other
// Amazon Web Services service, Amazon Simple Notification Service might place your
// account in the SMS sandbox. In [sandbox mode], you can send messages only to verified phone
// numbers. After you test your app while in the sandbox environment, you can move
// out of the sandbox and into production. For more information, see [SMS message settings for Amazon Cognito user pools]in the Amazon
// Cognito Developer Guide.
//
// Amazon Cognito evaluates Identity and Access Management (IAM) policies in
// requests for this API operation. For this operation, you must use IAM
// credentials to authorize requests, and you must grant yourself the corresponding
// IAM permission in a policy.
//
// # Learn more
//
// [Signing Amazon Web Services API Requests]
//
// [Using the Amazon Cognito user pools API and user pool endpoints]
//
// [SMS message settings for Amazon Cognito user pools]: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-sms-settings.html
// [Using the Amazon Cognito user pools API and user pool endpoints]: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pools-API-operations.html
// [sandbox mode]: https://docs.aws.amazon.com/sns/latest/dg/sns-sms-sandbox.html
// [Signing Amazon Web Services API Requests]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_aws-signing.html
// [Authentication flows]: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow-methods.html
// [Amazon Pinpoint]: https://console.aws.amazon.com/pinpoint/home/
func (c *Client) AdminInitiateAuth(ctx context.Context, params *AdminInitiateAuthInput, optFns ...func(*Options)) (*AdminInitiateAuthOutput, error) {
	if params == nil {
		params = &AdminInitiateAuthInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "AdminInitiateAuth", params, optFns, c.addOperationAdminInitiateAuthMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*AdminInitiateAuthOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Initiates the authorization request, as an administrator.
type AdminInitiateAuthInput struct {

	// The authentication flow that you want to initiate. Each AuthFlow has linked
	// AuthParameters that you must submit. The following are some example flows and
	// their parameters.
	//
	//   - USER_AUTH : Request a preferred authentication type or review available
	//   authentication types. From the offered authentication types, select one in a
	//   challenge response and then authenticate with that method in an additional
	//   challenge response.
	//
	//   - REFRESH_TOKEN_AUTH : Receive new ID and access tokens when you pass a
	//   REFRESH_TOKEN parameter with a valid refresh token as the value.
	//
	//   - USER_SRP_AUTH : Receive secure remote password (SRP) variables for the next
	//   challenge, PASSWORD_VERIFIER , when you pass USERNAME and SRP_A parameters..
	//
	//   - ADMIN_USER_PASSWORD_AUTH : Receive new tokens or the next challenge, for
	//   example SOFTWARE_TOKEN_MFA , when you pass USERNAME and PASSWORD parameters.
	//
	// All flows
	//
	// USER_AUTH The entry point for sign-in with passwords, one-time passwords, and
	// WebAuthN authenticators.
	//
	// USER_SRP_AUTH Username-password authentication with the Secure Remote Password
	// (SRP) protocol. For more information, see [Use SRP password verification in custom authentication flow].
	//
	// REFRESH_TOKEN_AUTH and REFRESH_TOKEN Provide a valid refresh token and receive
	// new ID and access tokens. For more information, see [Using the refresh token].
	//
	// CUSTOM_AUTH Custom authentication with Lambda triggers. For more information,
	// see [Custom authentication challenge Lambda triggers].
	//
	// ADMIN_USER_PASSWORD_AUTH Username-password authentication with the password
	// sent directly in the request. For more information, see [Admin authentication flow].
	//
	// USER_PASSWORD_AUTH is a flow type of [InitiateAuth] and isn't valid for AdminInitiateAuth.
	//
	// [Use SRP password verification in custom authentication flow]: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow.html#Using-SRP-password-verification-in-custom-authentication-flow
	// [InitiateAuth]: https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_InitiateAuth.html
	// [Using the refresh token]: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-the-refresh-token.html
	// [Admin authentication flow]: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow.html#Built-in-authentication-flow-and-challenges
	// [Custom authentication challenge Lambda triggers]: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-challenge.html
	//
	// This member is required.
	AuthFlow types.AuthFlowType

	// The ID of the app client where the user wants to sign in.
	//
	// This member is required.
	ClientId *string

	// The ID of the user pool where the user wants to sign in.
	//
	// This member is required.
	UserPoolId *string

	// The analytics metadata for collecting Amazon Pinpoint metrics.
	AnalyticsMetadata *types.AnalyticsMetadataType

	// The authentication parameters. These are inputs corresponding to the AuthFlow
	// that you're invoking. The required values depend on the value of AuthFlow :
	//
	//   - For USER_AUTH : USERNAME (required), PREFERRED_CHALLENGE . If you don't
	//   provide a value for PREFERRED_CHALLENGE , Amazon Cognito responds with the
	//   AvailableChallenges parameter that specifies the available sign-in methods.
	//
	//   - For USER_SRP_AUTH : USERNAME (required), SRP_A (required), SECRET_HASH
	//   (required if the app client is configured with a client secret), DEVICE_KEY .
	//
	//   - For ADMIN_USER_PASSWORD_AUTH : USERNAME (required), PASSWORD (required),
	//   SECRET_HASH (required if the app client is configured with a client secret),
	//   DEVICE_KEY .
	//
	//   - For REFRESH_TOKEN_AUTH/REFRESH_TOKEN : REFRESH_TOKEN (required), SECRET_HASH
	//   (required if the app client is configured with a client secret), DEVICE_KEY .
	//
	//   - For CUSTOM_AUTH : USERNAME (required), SECRET_HASH (if app client is
	//   configured with client secret), DEVICE_KEY . To start the authentication flow
	//   with password verification, include ChallengeName: SRP_A and SRP_A: (The
	//   SRP_A Value) .
	//
	// For more information about SECRET_HASH , see [Computing secret hash values]. For information about DEVICE_KEY
	// , see [Working with user devices in your user pool].
	//
	// [Computing secret hash values]: https://docs.aws.amazon.com/cognito/latest/developerguide/signing-up-users-in-your-app.html#cognito-user-pools-computing-secret-hash
	// [Working with user devices in your user pool]: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-device-tracking.html
	AuthParameters map[string]string

	// A map of custom key-value pairs that you can provide as input for certain
	// custom workflows that this action triggers.
	//
	// You create custom workflows by assigning Lambda functions to user pool
	// triggers. When you use the AdminInitiateAuth API action, Amazon Cognito invokes
	// the Lambda functions that are specified for various triggers. The ClientMetadata
	// value is passed as input to the functions for only the following triggers:
	//
	//   - Pre signup
	//
	//   - Pre authentication
	//
	//   - User migration
	//
	// When Amazon Cognito invokes the functions for these triggers, it passes a JSON
	// payload, which the function receives as input. This payload contains a
	// validationData attribute, which provides the data that you assigned to the
	// ClientMetadata parameter in your AdminInitiateAuth request. In your function
	// code in Lambda, you can process the validationData value to enhance your
	// workflow for your specific needs.
	//
	// When you use the AdminInitiateAuth API action, Amazon Cognito also invokes the
	// functions for the following triggers, but it doesn't provide the ClientMetadata
	// value as input:
	//
	//   - Post authentication
	//
	//   - Custom message
	//
	//   - Pre token generation
	//
	//   - Create auth challenge
	//
	//   - Define auth challenge
	//
	//   - Custom email sender
	//
	//   - Custom SMS sender
	//
	// For more information, see [Customizing user pool Workflows with Lambda Triggers] in the Amazon Cognito Developer Guide.
	//
	// When you use the ClientMetadata parameter, note that Amazon Cognito won't do
	// the following:
	//
	//   - Store the ClientMetadata value. This data is available only to Lambda
	//   triggers that are assigned to a user pool to support custom workflows. If your
	//   user pool configuration doesn't include triggers, the ClientMetadata parameter
	//   serves no purpose.
	//
	//   - Validate the ClientMetadata value.
	//
	//   - Encrypt the ClientMetadata value. Don't send sensitive information in this
	//   parameter.
	//
	// [Customizing user pool Workflows with Lambda Triggers]: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html
	ClientMetadata map[string]string

	// Contextual data about your user session, such as the device fingerprint, IP
	// address, or location. Amazon Cognito advanced security evaluates the risk of an
	// authentication event based on the context that your app generates and passes to
	// Amazon Cognito when it makes API requests.
	//
	// For more information, see [Collecting data for threat protection in applications].
	//
	// [Collecting data for threat protection in applications]: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-viewing-threat-protection-app.html
	ContextData *types.ContextDataType

	// The optional session ID from a ConfirmSignUp API request. You can sign in a
	// user directly from the sign-up process with an AuthFlow of USER_AUTH and
	// AuthParameters of EMAIL_OTP or SMS_OTP , depending on how your user pool sent
	// the confirmation-code message.
	Session *string

	noSmithyDocumentSerde
}

// Initiates the authentication response, as an administrator.
type AdminInitiateAuthOutput struct {

	// The outcome of successful authentication. This is only returned if the user
	// pool has no additional challenges to return. If Amazon Cognito returns another
	// challenge, the response includes ChallengeName , ChallengeParameters , and
	// Session so that your user can answer the challenge.
	AuthenticationResult *types.AuthenticationResultType

	// The name of the challenge that you're responding to with this call. This is
	// returned in the AdminInitiateAuth response if you must pass another challenge.
	//
	//   - WEB_AUTHN : Respond to the challenge with the results of a successful
	//   authentication with a passkey, or webauthN, factor. These are typically
	//   biometric devices or security keys.
	//
	//   - PASSWORD : Respond with USER_PASSWORD_AUTH parameters: USERNAME (required),
	//   PASSWORD (required), SECRET_HASH (required if the app client is configured
	//   with a client secret), DEVICE_KEY .
	//
	//   - PASSWORD_SRP : Respond with USER_SRP_AUTH parameters: USERNAME (required),
	//   SRP_A (required), SECRET_HASH (required if the app client is configured with a
	//   client secret), DEVICE_KEY .
	//
	//   - SELECT_CHALLENGE : Respond to the challenge with USERNAME and an ANSWER that
	//   matches one of the challenge types in the AvailableChallenges response
	//   parameter.
	//
	//   - MFA_SETUP : If MFA is required, users who don't have at least one of the MFA
	//   methods set up are presented with an MFA_SETUP challenge. The user must set up
	//   at least one MFA type to continue to authenticate.
	//
	//   - SELECT_MFA_TYPE : Selects the MFA type. Valid MFA options are SMS_MFA for
	//   SMS message MFA, EMAIL_OTP for email message MFA, and SOFTWARE_TOKEN_MFA for
	//   time-based one-time password (TOTP) software token MFA.
	//
	//   - SMS_MFA : Next challenge is to supply an SMS_MFA_CODE that your user pool
	//   delivered in an SMS message.
	//
	//   - EMAIL_OTP : Next challenge is to supply an EMAIL_OTP_CODE that your user
	//   pool delivered in an email message.
	//
	//   - PASSWORD_VERIFIER : Next challenge is to supply PASSWORD_CLAIM_SIGNATURE ,
	//   PASSWORD_CLAIM_SECRET_BLOCK , and TIMESTAMP after the client-side SRP
	//   calculations.
	//
	//   - CUSTOM_CHALLENGE : This is returned if your custom authentication flow
	//   determines that the user should pass another challenge before tokens are issued.
	//
	//   - DEVICE_SRP_AUTH : If device tracking was activated in your user pool and the
	//   previous challenges were passed, this challenge is returned so that Amazon
	//   Cognito can start tracking this device.
	//
	//   - DEVICE_PASSWORD_VERIFIER : Similar to PASSWORD_VERIFIER , but for devices
	//   only.
	//
	//   - ADMIN_NO_SRP_AUTH : This is returned if you must authenticate with USERNAME
	//   and PASSWORD directly. An app client must be enabled to use this flow.
	//
	//   - NEW_PASSWORD_REQUIRED : For users who are required to change their passwords
	//   after successful first login. Respond to this challenge with NEW_PASSWORD and
	//   any required attributes that Amazon Cognito returned in the requiredAttributes
	//   parameter. You can also set values for attributes that aren't required by your
	//   user pool and that your app client can write. For more information, see [AdminRespondToAuthChallenge].
	//
	// Amazon Cognito only returns this challenge for users who have temporary
	//   passwords. Because of this, and because in some cases you can create users who
	//   don't have values for required attributes, take care to collect and submit
	//   required-attribute values for all users who don't have passwords. You can create
	//   a user in the Amazon Cognito console without, for example, a required
	//   birthdate attribute. The API response from Amazon Cognito won't prompt you to
	//   submit a birthdate for the user if they don't have a password.
	//
	// In a NEW_PASSWORD_REQUIRED challenge response, you can't modify a required
	//   attribute that already has a value. In AdminRespondToAuthChallenge , set a
	//   value for any keys that Amazon Cognito returned in the requiredAttributes
	//   parameter, then use the AdminUpdateUserAttributes API operation to modify the
	//   value of any additional attributes.
	//
	//   - MFA_SETUP : For users who are required to set up an MFA factor before they
	//   can sign in. The MFA types activated for the user pool will be listed in the
	//   challenge parameters MFAS_CAN_SETUP value.
	//
	// To set up software token MFA, use the session returned here from InitiateAuth as
	//   an input to AssociateSoftwareToken , and use the session returned by
	//   VerifySoftwareToken as an input to RespondToAuthChallenge with challenge name
	//   MFA_SETUP to complete sign-in. To set up SMS MFA, users will need help from an
	//   administrator to add a phone number to their account and then call
	//   InitiateAuth again to restart sign-in.
	//
	// [AdminRespondToAuthChallenge]: https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminRespondToAuthChallenge.html
	ChallengeName types.ChallengeNameType

	// The challenge parameters. These are returned to you in the AdminInitiateAuth
	// response if you must pass another challenge. The responses in this parameter
	// should be used to compute inputs to the next call ( AdminRespondToAuthChallenge
	// ).
	//
	// All challenges require USERNAME and SECRET_HASH (if applicable).
	//
	// The value of the USER_ID_FOR_SRP attribute is the user's actual username, not
	// an alias (such as email address or phone number), even if you specified an alias
	// in your call to AdminInitiateAuth . This happens because, in the
	// AdminRespondToAuthChallenge API ChallengeResponses , the USERNAME attribute
	// can't be an alias.
	ChallengeParameters map[string]string

	// The session that must be passed to challenge-response requests. If an
	// AdminInitiateAuth or AdminRespondToAuthChallenge API request determines that
	// the caller must pass another challenge, Amazon Cognito returns a session ID and
	// the parameters of the next challenge. Pass this session Id in the Session
	// parameter of AdminRespondToAuthChallenge .
	Session *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationAdminInitiateAuthMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpAdminInitiateAuth{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpAdminInitiateAuth{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "AdminInitiateAuth"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpAdminInitiateAuthValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opAdminInitiateAuth(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opAdminInitiateAuth(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "AdminInitiateAuth",
	}
}