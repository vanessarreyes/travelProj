// TODO: find way to create a common js file for shared functions
// create card
function CreateCard(place, description, rowCount) {
    // create column
    var  col = document.createElement('div')
    col.setAttribute('class', 'col-sm')

    // create card div
    var card = document.createElement('div')
    card.setAttribute('class', 'card mx-auto h-100')

    // create card body div
    var cardBody = document.createElement('div')
    cardBody.setAttribute('class', 'card-body')

    // create card body content
    var cardTitle = document.createElement('h5')
    cardTitle.setAttribute('class', 'card-title')
    cardTitle.innerHTML = place

    var cardText = document.createElement('p')
    cardText.setAttribute('class', 'card-text')
    cardText.innerHTML = description

    var cardButton = document.createElement('a')
    cardButton.setAttribute('href', '#') // need to change this to link
    cardButton.setAttribute('class', 'btn btn-primary')
    cardButton.innerHTML = "View"

    // Add data to card body
    cardBody.appendChild(cardTitle)
    cardBody.appendChild(cardText)
    cardBody.appendChild(cardButton)

    // Add data to card
    card.appendChild(cardBody)

    // Add data to the column
    col.appendChild(card)

    // Add col to row
    var row = document.getElementById("cardRow" + rowCount)
    row.appendChild(col)
}

function EmptyDisplay(container) {
    // Create div for content
    div = document.createElement('div')
    div.setAttribute('id', 'emptyContent')

    // Create text
    text = document.createElement('h4')
    text.innerHTML = "Looks like you don't have any past travels. Let's start building your journal by putting your past travels!"
    
    // Create button
    button = document.createElement('a')
    button.setAttribute('href', 'http://localhost:400/PastTravelForm')
    button.setAttribute('role', 'button')
    button.setAttribute('class', 'btn btn-dark btn-lg')
    button.innerHTML = 'Get Started'
    
    // Add elements to div and container
    div.appendChild(text)
    div.appendChild(button)
    container.appendChild(div)

    // TODO: add some images to take up some space
}

// get cards for next travel ideas page
function LoadCards() {
    fetch("/PastTravelsCards")
    .then(response => response.text())
    .then(text => {
        var container = document.getElementById("cardContainer")
        try {
            const travelList = JSON.parse(text) // Try to parse the response as JSON
            if ((travelList).length == 0) {
                console.log("here")
                EmptyDisplay(container)
                return
            }
            var header = document.getElementById('headerContent')
            var button = document.createElement('a')
            button.setAttribute('href', 'http://localhost:400/PastTravelForm')
            button.setAttribute('role', 'button')
            button.setAttribute('class', 'btn btn-dark')
            button.innerHTML = "Add Past Travel"
            header.appendChild(button)

            // response was JSON
            var rowCount = 0
            var colCount = 0

            // Once we fetch the list, we iterate over it
            travelList.forEach(function(travel) {
                console.log(travel)
                if (colCount % 3 == 0) {
                    rowCount++
                    // add start row
                    var row = document.createElement('div')
                    row.setAttribute('class', 'row')
                    row.setAttribute('id', 'cardRow' + rowCount)
                    container.appendChild(row)
                }

                CreateCard(travel.place, travel.description, rowCount)

                // Add to column count
                colCount++
            })
        } catch(err) {
            // JSON was only one object
            const travel = JSON.parse(text) // Try to parse the response as JSON
            var row = document.createElement('div')
            row.setAttribute('class', 'row')
            row.setAttribute('id', 'cardRow' + rowCount)
            container.appendChild(row)
            CreateCard(travel.place, travel.description, rowCount)
        }
    })
}