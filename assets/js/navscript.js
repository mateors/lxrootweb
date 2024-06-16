let naccessName="client"; 
function arrayValueExist(findme, comparewithUrl){
  const urls=[];
  if(naccessName!="client"){
      urls['/orders']="/orders";

  }else if (naccessName=="client"){
      
      urls['/license']="/license";
      urls['/orders']="/orders";
      urls['/profile']="/profile";
      urls['/tickets']="/tickets";
      urls['/invoices']="/orders"
      urls['/security']="/profile"
      urls['/ticketnew']="/tickets"
      urls['/ticket/']="/tickets"
      urls['/paymethods']="/orders"
      urls['/activity']="/profile"

  }
  let val=urls[findme];
  if (val==comparewithUrl){
    return true;
  } else if (findme.includes(comparewithUrl)){
    return true;
  }
  return false;
}

$('#topMenu').findAll("li>a").each( (element, index) => {

  $(element).removeClass(["text-primary","font-semibold"]);
  if (arrayValueExist(location.pathname,element.pathname)==true){
      $(element).addClass(["text-primary","font-semibold"]);
  }

});