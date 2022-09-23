<script lang="ts">

  import { calculateFunding, count, isFetching, requestFtx, setNewCount } from './lib/store.js';

  async function onSubmit(e) {
    const formData = new FormData(e.target);

    const data = {};
    for (let field of formData) {
      const [key, value] = field;
      data[key] = value;
    }

    const stuff = await requestFtx(data["perp"], data["trade-start"], data["trade-end"])
    let humanRead = (num,intSep = ',',floatSep = '.') => {
      return new Intl
          .NumberFormat('en-US')
          .format(num)
          .replaceAll('.',floatSep)
          .replaceAll(',',intSep);
      }
    setNewCount(humanRead(calculateFunding(stuff.result, data["size"], data["side"])))
  }

</script>

<main>
<form on:submit|preventDefault={onSubmit}>
    <div>
        <label for="name">Notional size in $</label>
        <input
          type="number"
          required="true"
          id="size"
          name="size"
          value=""
        />
    </div>
    <div>
        <label for="name">Futures name</label>
        <input
          type="text"
          id="perp"
          required="true"
          name="perp"
          value="BTC-PERP"
        />
    </div>
    <div>
      <!-- add label and select for the new field with options: long or short -->
      <label for="name">Long/Short</label>
      <select id="side" name="side">
        <option value="long">Long</option>
        <option value="short">Short</option>
      </select>


    </div>
    <div>
      <!-- start date selector -->
      <label for="name">Start date</label>
      <input
        type="date"
        id="start"
        name="trade-start"
        required="true"
      />
  </div>
  <div>
    <!-- end date selector -->
    <label for="name">End date</label>
    <input
      type="date"
      id="end"
      name="trade-end"
      required="true"
    />
  </div>

  <button type="submit">
	Calculate funding rate {#if $isFetching}🌀{/if}
</button>
  </form>

<pre>
	Calculated funding payments is {$count}$
</pre>
</main>
