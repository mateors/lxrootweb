var navItems = document.querySelectorAll("nav .nav-item");

navItems.forEach(navItem => {
    if(navItem) {
        navItem.onclick = () => {
            navItems.forEach(navItem => {
                navItem.classList.remove("active");
            })
            navItem.classList.add("active");
        }
    }
})

/* changing sidebar active according to page */
var navigations = document.querySelectorAll("aside nav > a.nav-item");

navigations.forEach(navigation => {
    // var currentPageLocation = location.href;
    var locationUrlSplit = location.pathname.split("/");
    var rootPage = locationUrlSplit[1];
    
    var navUrlSplit = navigation.pathname.split("/");
    var navUrl = navUrlSplit[1];

    if(navUrl === rootPage){
        navigation.classList.add("active")
    }
});

/* change tab active according to page */
var tabButtons = document.querySelectorAll(".tab-bar a.tab-button");
        
tabButtons.forEach(tabButton => {
    // var currentPageLocation = location.href;
    var locationUrlSplit = location.pathname.split("/");
    var subPage = locationUrlSplit[2];
    // console.log(subPage)
    var tabUrlSplit = tabButton.pathname.split("/");
    var tabUrl = tabUrlSplit[2];
    console.log(tabUrl)

    if(tabUrl === subPage){
        tabButton.classList.add("active-tab")
    }
});

/* dropdown functionality */
document.addEventListener("click", e => {

    //console.log(e.target,e.target.closest("[data-dropdown]"));
    const isDropdownBtn = e.target.matches("[data-dropdown-btn]");
    if(!isDropdownBtn && e.target.closest("[data-dropdown]") != null){
        return;
    }

    let currentDropdown
    if(isDropdownBtn) {
        currentDropdown = e.target.closest("[data-dropdown]");
        currentDropdown.classList.toggle("menu-open");
    }

    document.querySelectorAll("[data-dropdown].menu-open").forEach(dropdown => {
        var clsDropdownBtn = dropdown.querySelector("[data-dropdown-cls]");
        if(clsDropdownBtn) {
            clsDropdownBtn.onclick = () => {
                dropdown.classList.remove("menu-open");
            }
        }
        if(dropdown === currentDropdown){
            return;
        }
        dropdown.classList.remove("menu-open");
        
    })
});

/* tabs functionality */
var tabs = document.querySelector(".tab-bar");
var tabButton = document.querySelectorAll(".tab-bar > .tab-button");
var contents = document.querySelectorAll(".tab-content-wrap > .tab-content");

if (tabs) {
    tabs.onclick = e => {
        var id = e.target.dataset.id;
        console.log(id);
        if (id) {
            tabButton.forEach(btn => {
                btn.classList.remove("active-tab");
            });
            e.target.classList.add("active-tab");

            contents.forEach(content => {
                content.classList.remove("active-content");
            });
            var element = document.getElementById(id);
            element.classList.add("active-content");
        }
    }
}


var innerTabWrapper = document.querySelectorAll(".tab-content-wrap > .tab-content")
innerTabWrapper.forEach(innerTabWrap => {
    var innerTabs = innerTabWrap.querySelector(".inner-tab-bar");
    var innerTabButton = innerTabWrap.querySelectorAll(".inner-tab-button");
    var innerContents = innerTabWrap.querySelectorAll(".inner-tab-content-wrap .tab-content");

    if (innerTabs) {
        innerTabs.onclick = e => {
            var id = e.target.dataset.id;
            if (id) {
                innerTabButton.forEach(btn => {
                    btn.classList.remove("active-tab");
                });
                e.target.classList.add("active-tab");
                

                innerContents.forEach(content => {
                    content.classList.remove("active-content");
                });
                var element = document.getElementById(id);
                element.classList.add("active-content");
            }
        }
    }
    
})

var modalTabs = document.querySelectorAll(".modal-tab-bar");
var modalTabButton = document.querySelectorAll(".modal-tab-bar .tab-button");
var modalTabContents = document.querySelectorAll(".modal-tab-content-wrap > .modal-tab-content");

modalTabs.forEach(modalTab => {
    if (modalTab) {
        modalTab.onclick = e => {
            var id = e.target.dataset.id;
            console.log(id);
            if (id) {
                modalTabButton.forEach(btn => {
                    btn.classList.remove("active-tab");
                });
                e.target.classList.add("active-tab");
    
                modalTabContents.forEach(content => {
                    content.classList.remove("active-content");
                });
                var element = document.getElementById(id);
                element.classList.add("active-content");
    
            }
        }
    }
})

/* modal */
var modalButtons = document.querySelectorAll("[data-modal]");
// var modalContents = document.querySelectorAll(".modal-wrap");

modalButtons.forEach(modalButton => {
    if(modalButton) {
        modalButton.onclick = e => {
            var modalId = e.target.dataset.modal;
            console.log(modalId);

            var modalContent = document.getElementById(modalId);
            modalContent.classList.add("active-modal");
            document.querySelector("body").classList.add("modal-open");

            var closeModal = modalContent.querySelectorAll("[data-modal-closer]");
            closeModal.forEach(clsModel => {
                if(clsModel) {
                    clsModel.onclick = () => {
                        modalContent.classList.remove("active-modal");
                        document.querySelector("body").classList.remove("modal-open");
                    }
                }
            })
            
            
            var doneModal = modalContent.querySelector("[data-modal-done]");
            if(doneModal) {
                doneModal.onclick = () => {
                    modalContent.classList.remove("active-modal");
                    document.querySelector("body").classList.remove("modal-open");

                    var toast = document.querySelector("[data-toast]");
                    toast.classList.add("active");
                    setTimeout(() => {
                        toast.classList.remove("active");
                    }, 3000)
                }
            }


        }
    }
})

/* collapsible */
var coll = document.getElementsByClassName("collapsible");
var i;

for (i = 0; i < coll.length; i++) {
    
    coll[i].addEventListener("click", function() {

        this.classList.toggle("active");
        var content = this.nextElementSibling;

        if (content.style.maxHeight){
            content.style.maxHeight = null;
        } else {
            content.style.maxHeight = content.scrollHeight + "px";    
        } 
    });
};

/* message */
var message = document.querySelector("[data-message]");

if(message) {
    var clsMessage = message.querySelector("[data-close-message]");
    clsMessage.onclick = () => {
        message.classList.add("hide");
    }
}

/* toast */
var toast = document.querySelector("[data-toast]");
var clsToast = document.querySelector("[data-toast-cls]");

function openToast(){
    toast.classList.add("active");

    setTimeout(() => {
        toast.classList.remove("active");
    }, 3000)
}

if(toast){
    clsToast.onclick = () => {
        toast.classList.remove("active");
    }
}

/* loading screen */
function formSubmit(){
    var loading = document.getElementById("loading-screen");
    loading.classList.add("active");

    setTimeout(() => {
        loading.classList.remove("active");
    }, 3000)
}


/* copy clipboard */
document.addEventListener("click", e => {
    var isCopyBtn = e.target.matches("[data-copy-button]");
    if(!isCopyBtn && e.target.closest("[data-copy-clipboard]") != null){
        return;
    }

    let currentClipboard
    if(isCopyBtn) {
        currentClipboard = e.target.closest("[data-copy-clipboard]");
        var copyContent = currentClipboard.querySelector(".copy-content");

        // if (document.selection){
        //     var div = currentClipboard.createTextRange();

        //     div.moveToElementText(copyContent);
        //     copyContent.select();
        // } else {
        //     var div = document.createRange();

        //     div.setStartBefore(copyContent);
        //     div.setEndAfter(copyContent);

        //     window.getSelection().addRange(div);
        // }
        copyContent.select();    
        document.execCommand("Copy");
    }

    document.querySelectorAll("[data-copy-clipboard]").forEach(clipboard => {
        if(clipboard === currentClipboard){
            return;
        }
        clipboard.classList.remove("menu-open");
    })

    
});


/* custom range */
const rangeInputs = document.querySelectorAll('input[type="range"]');

function handleInputChange(e) {
    let target = e.target
    if (e.target.type !== 'range') {
        target = document.getElementById('range')
        console.log(target)
    } 
    const min = target.min
    const max = target.max
    const val = target.value
    target.style.backgroundSize = (val - min) * 100 / (max - min) + '% 100%'
}

rangeInputs.forEach(input => {
    input.addEventListener('input', handleInputChange)
});


/* password toggle */
document.addEventListener("click", e => {
    let togglePassIcon = e.target.matches("[data-toggle-passicon]");
    if( !togglePassIcon && e.target.closest("[data-toggle]") != null){
        return;
    }

    let currentToggleField;
    let currentPasswordField;
    if(togglePassIcon){
        currentToggleField = e.target.closest("[data-toggle]");
        currentToggleField.classList.toggle("active");
        // currentToggleField.querySelector("[data-toggle-passicon]").innerHTML = "visibility_off";

        currentPasswordField = currentToggleField.querySelector("input");
        if(currentPasswordField.type == "password"){
            currentPasswordField.type = "text";
        } else {
            currentPasswordField.type = "password";
        }
        
        
    }


})

// header scroll effect
window.addEventListener('scroll', function() {
    const header = document.getElementById('main-header');
    if (!header) return; // If #main-header doesn't exist, exit the function

    if (window.scrollY > 100) {
        header.classList.remove('bg-transparent');
        header.classList.remove('border-b-[#808afd]');
        header.classList.add('bg-white');
        header.classList.add('border-b-slate-200');
        header.querySelectorAll('.nav-item').forEach(function(navItem) {
            navItem.classList.remove('text-white');
            navItem.classList.add('text-primaryDeep');
        });

        header.querySelectorAll('.logo-item > svg path').forEach(svg => {
            svg.style.fill = '#5865F2'
        });

        header.querySelectorAll('.nav-item > svg path').forEach(svg => {
            svg.style.stroke = '#515c79d1'
        });
    } else {
        header.classList.remove('bg-white');
        header.classList.remove('border-b-slate-200');
        header.classList.add('bg-transparent');
        header.classList.add('border-b-[#808afd]');
        header.querySelectorAll('.nav-item').forEach(function(navItem) {
            navItem.classList.add('text-white');
            navItem.classList.remove('text-primaryDeep');
        });

        header.querySelectorAll('.logo-item > svg path').forEach(svg => {
            svg.style.fill = '#fff'
        })

        header.querySelectorAll('.nav-item > svg path').forEach(svg => {
            svg.style.stroke = '#edededd1'
        });
    }
});

// Select the message container
const messageContainer = document.querySelector('[data-message-container]');

// Function to create a new toast message
function createToastMessage(message) {

    const toastMessage = document.createElement('div');
    toastMessage.classList.add('w-fit', 'message', 'px-5', 'py-3', 'mb-5', 'bg-secondary', 'text-white', 'flex', 'gap-5', 'items-center', 'justify-between', 'transition');
    toastMessage.dataset.toastmessage = true;

    // Set the content of the toast message
    toastMessage.innerHTML = `
        <p class="text-sm inline">${message}</p>
        <button data-close-message class="text-left sm:float-right text-red-500">
            DISMISS
        </button>
    `;
    // Add event listener to the close button
    const closeButton = toastMessage.querySelector('[data-close-message]');
    closeButton.addEventListener('click', () => {
        toastMessage.remove();
    });
    return toastMessage;
}

// Function to show the toast message
function showToastMessage(message) {

    const toastMessage = createToastMessage(message);
    messageContainer.append(toastMessage);

    // Shift existing messages upwards
    const messages = messageContainer.querySelectorAll('[data-toastmessage]');
    messages.forEach((msg, index) => {

        msg.classList.add('slide-up');
        msg.style.bottom = `${(index + 1) * 200}px`; // Adjust as per your layout

        setTimeout(() => { msg.classList.add("opacity-0"); }, 5000);
        setTimeout(() => { msg.classList.add("hidden");    }, 6000);
    });
    
}

// Example usage
const toastButtons = document.querySelectorAll('[data-toast-button]');
toastButtons.forEach(button => {
    button.addEventListener('click', () => {
        const message = button.dataset.toastmessage;
        console.log(message);
        showToastMessage(message);
    });
});


function elmDivSpinner(isLarger){
    //console.log(typeof isLarger, typeof isLarger === 'undefined');
    let divspn=document.createElement("div");
    if (typeof isLarger === 'undefined'){
        divspn.className="spinner";
    }else{
        divspn.classList.add("spinner","s32");
    }
    return divspn;
}

//showToastMessage("Congratulations! You have successfully registered with LxRoot. 🎉");
let validator = (selector, msg) => {

  let fVal=$(selector).val();
  let errMsg=$(selector).parent().next().text();
  let count=$(selector).parent().childrenLength();
  if (msg==""){ msg=errMsg; }

  if(fVal.trim()===""){
    $(selector).addClass('error');
    $(selector).parent().next().removeClass("hidden");
    $(selector).parent().childItem(count-1).removeClass("hidden");

  }else{
    $(selector).removeClass('error');
    $(selector).parent().next().addClass("hidden");
    $(selector).parent().childItem(count-1).addClass("hidden");
  }
  $(selector).parent().next().text(msg);
}