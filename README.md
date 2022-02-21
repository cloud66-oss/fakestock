<img src="http://cdn2-cloud66-com.s3.amazonaws.com/images/oss-sponsorship.png" width=150/>

# FakeStock

This is a simple API for test stock prices. 

*NOTE: All data produced by this service is randomly generated and is purely for test and demo purposes.*

The API has a few endpoints:

- `/_ping` returns OK 200 with the build commit and the duration since the service has been up.
- `/exchanges` returns the exchanges supported by the service.
- `/tickers` returns all tickers stored in the system.
- `/tickers/:symbol` returns info on the given symbol.

All symbols are loaded from `tickers.csv` when the service starts and allocated a random price between 1 and 1000. The service then randomly changes the price every 5 seconds. 1/3 of the tickers will not change while the other thirds will either go up or down by 10 units. 

## Hosting
A hosted version of this service runs on https://test.marketinfo.dev/ and is provided by Cloud 66 as a sample. You can host the service inside or outside a Docker container if you wish. 

## Credits
`tickers.csv` is downloaded from https://dumbstockapi.com/

This repository is maintained by Cloud 66 (https://www.cloud66.com)

