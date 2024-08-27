# Altinn API Example
This application calls a stateless Altinn skjema and gets the skjema results in response back.

Cross Origin requests are allowed here.

To find the relevant endpoints for you skjema, simply look at the network traffick tool in your browser when you use a Altinn skjema from the browser. We can use the same API as the browser does to fetch data

This application is the backend of a complete example on how a website can use Altinn skjema via its API endpoints. The frontend application is located [here](https://github.com/brreg/idporten-vite-spa-example)

### How it works

1. The user calls `/company` to get the available companies that can use this skjema, the request must provide a valid ID-porten access-token in the header.
2. ID-porten access-token is exchanged for a Altinn platform token at `https://platform.tt02.altinn.no/authentication/api/v1/exchange/id-porten`.
3. We find the available companies for the logged in user, using Altinn platform token in the header `https://brg.apps.tt02.altinn.no/brg/lpid-wallet-2024/api/v1/parties?allowedtoinstantiatefilter=true`.
4. `/company` returns a list with companies.
5. The user calls `/lpid` to get a credential offering, providing ID-porten access-token and PartyID in the header (PartyID comes from the list returned by `/company`).
6. ID-porten access-token is exchanged for a Altinn platform token at `https://platform.tt02.altinn.no/authentication/api/v1/exchange/id-porten`.
7. The Altinn skjema API endpoint is called with Altinn plattform token and PartyID in the header `https://brg.apps.tt02.altinn.no/brg/lpid-wallet-2024/v1/data?dataType=model&includeRowId=true&language=nb`.
8. `/lpid` returns the response, containing a LPID credential offering.

### Important Authentication Information

When integrating with ID-porten and Altinn tokens, it is common to encounter a fetch operation error. This usually happens under the following scenario:

1. **Random Daily Manager Selection:** The application retrieves a random daily manager through ID-porten, which is a necessary step for demonstrating varied user scenarios.
2. **Altinn Login Requirement:** If the randomly selected manager has never logged into Altinn, the token exchange will fail, preventing further actions within Altinn.

**Solution:** To resolve this issue, manually log in with the affected user account at least once at [Altinn's test environment](https://info.tt02.altinn.no/). This pre-authentication step ensures that the user's credentials are recognized by Altinn, facilitating successful token exchanges in subsequent automated processes.

Please ensure this step is followed during initial setup or when testing with new random users to avoid disruptions in service.
