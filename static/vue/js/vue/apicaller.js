
function handleErrorMsg(that, error) {
  const splitLines = str => str.split(/\r?\n/);

  console.log("Error is ", error)
  let err = error.bodyText
  if (err == '') {
    err = 'Unknown error ' + `$error`
  }
  that.$store.commit('clearMsgText') 
  that.$store.commit('resDatalog', splitLines(err))
  that.$store.commit('errorText', err)
}

export default {
  CallDataService(that, req) {
    return that.$http.post("CallDataService", JSON.stringify(req), { headers: { "content-type": "application/json" } })
  },
  ApproveCmt(that, params, Ok) {
    let req = { method: 'ApproveCmt', Params: params }
    this.CallDataService(that, req).then(result => {
      console.log('Call terminated ', result.data)
      that.$store.commit('clearMsgText') 
      that.$store.commit('resDatalog', [result.data.Status])
    }, error => {
      handleErrorMsg(that, error)
    });
  },
  DoLogin(that, params, Ok) {
    let req = { method: 'DoLogin', Params: params }
    this.CallDataService(that, req).then(result => {
      console.log('Call terminated ', result.data)
      that.$store.commit('clearMsgText') 
      that.$store.commit('resDatalog', [result.data.Status])
      that.$store.commit('storeToken', result.data.Token) 
    }, error => {
      handleErrorMsg(that, error)
    });
  }
}