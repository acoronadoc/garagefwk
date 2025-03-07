/* Init app */
window.addEventListener("load", function( event ) {
    execute( "/api/admin?url=" + encodeURIComponent( window.location.pathname ) );
});

/* Navigate */
window.addEventListener('popstate', function (event) {
    var stateObj = event.state;

    navigate( window.location.href );
});

window.addEventListener("click", (event) => {
    event.preventDefault();

    var link = getLink( event.target );
    if ( link ) {
        navigate( link.getAttribute("href") );
        return false;
    }
    
});

function navigate( url ) {
    execute( "/api/admin?url=" + encodeURIComponent( url ) );
    history.pushState( { }, "", url );

    document.querySelector("#header-menu").classList.remove('open');
}

function getLink( target ) {
    
    while ( target.tagName != "A" ) {
        target = target.parentNode;

        if ( !target ) return null;
    }

    if ( target.href.endsWith( "#" ) ) return null;

    return target;
}

/* Requests */
function execute( url, reg, onSuccess ) {
    var post=JSON.stringify(reg);

    fetch( url, {
          method: 'POST',
          body: post,
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
            },
        } ).then(function(response) {
            response.json().then(function(p) {
                if ( !onSuccess )
                    render( p );
                else
                    onSuccess( p );
            });
    });
}

function render( p ) {
    if ( p.parts )
        for ( var i=0; i<p.parts.length; i++ ) {

            if ( partScripts[ p.parts[i].component ] )
                partScripts[ p.parts[i].component ]( p.parts[i] );

            }
    }

var partScripts={};

partScripts.clean=function(p) {
    document.querySelector( "#appcontent" ).innerHTML = "";
}    

partScripts.eval=function(p) {
    eval( p.script );
} 

partScripts.renderTag=function(p) {
    tag = document.createElement( p.tag );

    Object.keys( p ).forEach( k=> {
        if ( !isInternalTag(k) ) {
            var s = p[ k ]
            if ( typeof s === 'object' ) s = JSON.stringify( s )

            if ( k == "innerHTML" ) tag.innerHTML = s
            else tag.setAttribute( k, s );
        }
    });

    if ( p.selector ) {
        var el = document.querySelector( p.selector );
        el.appendChild( tag );
    } else if ( p.insertBefore ) {
        var el = document.querySelector( p.insertBefore );
        el.parentNode.insertBefore( tag, el );
    } else {
        document.querySelector( "#appcontent" ).appendChild( tag );
    }
}    

partScripts.updateTag=function(p) {
    tag = document.querySelector( p.selector );

    Object.keys( p ).forEach( k=> {
        if ( !isInternalTag(k) ) {
            var s = p[ k ]
            if ( typeof s === 'object' ) s = JSON.stringify( s )

            tag.setAttribute( k, s );
        }
    });
}    

partScripts.deleteTag=function(p) {
    tag = document.querySelector( p.selector );

    document.querySelectorAll( p.selector ).forEach(e => e.remove());
}    

function isInternalTag( t ) {
    if ( t == "tag" || 
        t == "component" || 
        t == "selector" ||
        t == "updateTag" ||
        t == "renderTag" ||
        t == "insertBefore" 
     ) return true
}

/* Forms */
function saveForm(event) {
    event.preventDefault();
    form = getParentElement( event.target, "FORM" )
    values = {}

    form.querySelectorAll("input").forEach( function(el) { values[ el.getAttribute("name") ] = el.value; } );
    form.querySelectorAll("select").forEach( function(el) { values[ el.getAttribute("name") ] = el.value; } );
    form.querySelectorAll("textarea").forEach( function(el) { values[ el.getAttribute("name") ] = el.value; } );

    execute( "/api/admin?url=" + encodeURIComponent( window.location.pathname ), { "saveForm": values } );

    return false;
}

function checkForm(event) {
    values = { "checkForm": {} }
    
    form = document.querySelector( "#appcontent form" )

    if ( event ) {
        event.preventDefault();
        form = getParentElement( event.target, "FORM" )
        values["field"] = event.target.name
    }

    form.querySelectorAll("input").forEach( function(el) { values["checkForm"][ el.getAttribute("name") ] = el.value; } );
    form.querySelectorAll("select").forEach( function(el) { values["checkForm"][ el.getAttribute("name") ] = el.value; } );
    form.querySelectorAll("textarea").forEach( function(el) { values["checkForm"][ el.getAttribute("name") ] = el.value; } );

    execute( "/api/admin?url=" + encodeURIComponent( window.location.pathname ), values );

    return false;
}

/* Utils */
function getParentElement( el, parentTagName ) {
    form = el;

    for ( var i=0; i<20; i++ ) {
        if ( form.tagName == parentTagName )
            return form;
        
        form = form.parentNode
    }

    return null;
}