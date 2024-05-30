# Altinn API Example
This application calls a stateless Altinn skjema and gets the skjema results in response back.

Cross Origin requests are allowed here.

To find the relevant endpoints for you skjema, simply look at the network traffick tool in your browser when you use a Altinn skjema from the browser. We can use the same API as the browser does to fetch data

### How it works
1. The user calls `/company` to get the available companies that can use this skjema, the request must provide a valid ID-porten access-token in the header.
2. ID-porten access-token is exchanged for a Altinn platform token at `https://platform.tt02.altinn.no/authentication/api/v1/exchange/id-porten`.
3. We find the available companies for the logged in user, using Altinn platform token in the header `https://brg.apps.tt02.altinn.no/brg/lpid-wallet-2024/api/v1/parties?allowedtoinstantiatefilter=true`.
4. `/company` returns a list with companies.
5. The user calls `/lpid` to get a credential offering, providing ID-porten access-token and PartyID in the header (PartyID comes from the list returned by `/company`).
6. ID-porten access-token is exchanged for a Altinn platform token at `https://platform.tt02.altinn.no/authentication/api/v1/exchange/id-porten`.
7. The Altinn skjema API endpoint is called with Altinn plattform token and PartyID in the header `https://brg.apps.tt02.altinn.no/brg/lpid-wallet-2024/v1/data?dataType=model&includeRowId=true&language=nb`.
8. `/lpid` returns the response, containing a LPID credential offering.