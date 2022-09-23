import { writable } from "svelte/store";

export const count = writable(0);
export const isFetching = writable(false);

export function setNewCount(fundingPayment) {
  isFetching.set(true);
  return new Promise((res) =>
    setTimeout(() => {
      res(count.set(fundingPayment));
      isFetching.set(false);
    }, 1000)
  );
}

export async function requestFtx(perp, startTime, endTime) {
  // fetch https://ftx.com/api/funding_rates with as parameter future, startTime, endTime
  // return the result
  // https://ftx.com/api/funding_rates?future=BTC-PERP&start_time=1610000000&end_time=1610000000
  const url = `http://api.bnkd.me/fundings?future=${perp}&start_time=${startTime}&end_time=${endTime}`;
  const response = await fetch(url, {
    method: "GET",
    redirect: "follow",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return response.json();
}

export function calculateFunding(fundings, positionSize, side) {
  let result = 0;
  fundings.forEach((f) => {
    console.log(f);
    // if side is long, then we need to pay the funding rate
    // if side is short, then we need to receive the funding rate
    if (side === "long") {
      // if rate is positive, then we need to pay the rate
      // if rate is negative, then we need to receive the rate
      if (f.rate < 0) {
        result -= f.rate * positionSize;
      } else {
        result += f.rate * positionSize;
      }
    } else if (side === "short") {
      // if rate is positive, then we need to receive the rate
      // if rate is negative, then we need to pay the rate
      if (f.rate < 0) {
        result += f.rate * positionSize;
      } else {
        result -= f.rate * positionSize;
      }
    }
  });
  return result;
}

function toTimestamp(strDate) {
  var datum = Date.parse(strDate);
  return datum / 1000;
}
