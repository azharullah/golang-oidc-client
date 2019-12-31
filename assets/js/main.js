$(document).ready(function(){
  
  $("#verifyToken").on("click", function() {
    $.post(
      "/oidc/verify",
      function(data) {
        if(data.Success == true) {
          alert("Verification successful! Received data: " + JSON.stringify(JSON.parse(data.Msg), "", 2))
        }
        else {
          alert("Failed to verify token, error: " + data.Msg);
        }
      }
  );
  })

});