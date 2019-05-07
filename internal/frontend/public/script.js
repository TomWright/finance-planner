function newPage(id) {
    return {
        ID: id,
        Show: function () {
            console.log("Showing page: " + this.ID);
            let x = document.querySelectorAll('#' + this.ID + '.page');
            if (typeof x != 'undefined' && x.length > 0) {
                x.forEach(function (p) {
                    if (!p.classList.contains('active')) {
                        p.classList.add('active')
                    }
                })
            }
        },
        Hide: function () {
            console.log("Hiding page: " + this.ID);
            let x = document.querySelectorAll('#' + this.ID + '.page');
            if (typeof x != 'undefined' && x.length > 0) {
                x.forEach(function (p) {
                    if (p.classList.contains('active')) {
                        p.classList.remove('active')
                    }
                })
            }
        },
    };
}

function showPage(id) {
    let want = "page-" + id;
    for (let i in pages) {
        if (pages[i].ID === want) {
            pages[i].Show();
        } else {
            pages[i].Hide();
        }
    }
}

function getTransactions(profileName, callback) {
    callback = callback || function () {
    };

    let xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function () {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(JSON.parse(xmlHttp.responseText));
    };
    xmlHttp.open("GET", backendBaseUrl + '/' + profileName + '/transactions', true); // true for asynchronous
    xmlHttp.send(null);
}

function getStats(profileName, callback) {
    callback = callback || function () {
    };

    let xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function () {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(JSON.parse(xmlHttp.responseText));
    };
    xmlHttp.open("GET", backendBaseUrl + '/' + profileName + '/transactions/stats', true); // true for asynchronous
    xmlHttp.send(null);
}

function submitProfileLoginForm(e) {
    if (e.preventDefault) e.preventDefault();
    let profileInput = document.querySelector('#page-login form input[name=profile]');
    let profileName = profileInput.value;
    if (profileName === '') {
        alert("Please enter a profile name");
        return false;
    }

    refreshProfileOverview(profileName, function () {
        showPage('profile-overview');
    });
}

function refreshProfileOverview(profileName, callback) {
    callback = callback || function(){};
    let transactionsLoaded = false;
    let statsLoaded = false;
    let showPageIfPossible = function () {
        if (transactionsLoaded && statsLoaded) {
            callback();
        }
    };
    getTransactions(profileName, function (transactions) {
        let x = document.querySelectorAll('#page-profile-overview div.transactions ul');
        if (typeof x == 'undefined') {
            transactionsLoaded = true;
            return
        }
        x.forEach(function (ts) {
            ts.innerHTML = "";
        });
        let transactionsIn = document.querySelector('#page-profile-overview div.transactions.in ul');
        let transactionsOut = document.querySelector('#page-profile-overview div.transactions.out ul');

        transactions.data.forEach(function (t) {
            let tagStr = "";
            if (typeof t.tags != 'undefined' && t.tags.length > 0) {
                tagStr = ' [' + t.tags.join(', ') + ']'
            }
            let li = document.createElement("li");
            li.innerText = t.label + ": " + Math.abs(t.amount) + tagStr;
            if (t.amount > 0) {
                transactionsIn.appendChild(li)
            } else {
                transactionsOut.appendChild(li)
            }
        });

        transactionsLoaded = true;
        showPageIfPossible();
    });
    getStats(profileName, function (stats) {
        let x = document.querySelector('#page-profile-overview div.stats');
        if (typeof x == 'undefined') {
            statsLoaded = true;
            return
        }
        x.innerHTML = "";

        let viewableStats = [
            {
                label: "Sum",
                value: stats.sum
            }
        ];

        viewableStats.forEach(function (v) {
            let li = document.createElement("li");
            li.innerText = v.label + ": " + v.value;
            x.appendChild(li)
        });

        statsLoaded = true;
        showPageIfPossible();
    });
}

window.onload = function () {
    window.pages = {
        profileInput: newPage("page-login"),
        profileOverview: newPage("page-profile-overview")
    };

    let loginForm = document.querySelector('#page-login form');
    if (loginForm.attachEvent) {
        loginForm.attachEvent("submit", submitProfileLoginForm);
    } else {
        loginForm.addEventListener("submit", submitProfileLoginForm);
    }

    showPage("login")
};