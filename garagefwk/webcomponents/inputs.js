
class AppInputElement extends HTMLElement {

    static observedAttributes = [ "name", "label", "type", "options",  "value", "onchange" ];

    elementAdded = false;
  
    constructor() { 
        super(); 
    }
  
    connectedCallback() {
      //console.log("Custom element added to page.");
      this.elementAdded = true;
      this.render();
    }
   
    attributeChangedCallback(name, oldValue, newValue) {
      //console.log(`Attribute ${name} has changed . ${oldValue} => ${newValue}`);
      if ( name == "value" && this.elementAdded ) {
        this.setValue( newValue );
        return;
      }

      if ( this.elementAdded ) 
        this.render();
    }

    stateEvent() {
        var cl="";
        var el = this.querySelector("input, select, textarea");

        if ( el.value != "" ) cl = "full";

        if ( el.tagName == "SELECT" && el.selectedIndex != -1 &&
          el.options[el.selectedIndex].innerHTML.trim() != ""
        ) {
            cl = "full";
        }

        if ( document.activeElement === el ) {
          cl = "full focus";
        } 

        this.querySelector(".app-textbox-wrapper").className = "app-textbox-wrapper " + cl;
    }

    setValue( value ) {
      var el = this.querySelector("input, select, textarea");

      value = value.replace( "\\n", "\n" );

      el.value = value;
    }

    render() {
        var el = this;
        var editor = null;

        var ttype = this.getAttribute( "type" );
        var value = this.getAttribute( "value" );
        var onchange = this.getAttribute( "onchange" );
        if ( !ttype ) ttype = "textbox";
        if ( !value ) value = "";
        if ( !onchange ) onchange = "";

        this.innerHTML = "<div class='app-textbox-wrapper'></div>";

        this.querySelector("div.app-textbox-wrapper").appendChild( document.createRange().createContextualFragment( "<label>"+this.getAttribute( "label" )+"</label>" ) );

        if ( ttype == "textbox" ) {
          editor = this.createElement( "input", { 'name': this.getAttribute( "name" ) } )
        }

        if ( ttype == "readonly" ) {
          editor = this.createElement( "input", { 'name': this.getAttribute( "name"), 'readonly': 'readonly' } )
        }

        if ( ttype == "select" ) {
          editor = this.createElement( "select", { 'name': this.getAttribute( "name" ) } );
          this.fillSelect( editor, JSON.parse( this.getAttribute( "options" ) ) );
        }

        if ( ttype == "textarea" ) {
          editor = this.createElement( "textarea", { 'name': this.getAttribute( "name" ) } )
        }

        this.querySelector("div.app-textbox-wrapper").append( editor );
        this.querySelector("label").addEventListener('click', function (event) { editor.focus(); });

        editor.addEventListener('focus', function (event) { el.stateEvent(); });
        editor.addEventListener('blur', function (event) { el.stateEvent(); });

        /*if ( onchange != "" ) {
          editor.addEventListener('change', function (event) { onchange=onchange.replace("javascript:",""); eval( onchange ); });
        }*/

        this.setValue( value );
        el.stateEvent();
    }

    createElement( tagname, params, innerHTML ) {
      var el = document.createElement( tagname );

      var paramKeys = Object.keys( params );

      for ( var i=0; i<paramKeys.length;  i++ ) {
        el.setAttribute( paramKeys[i], params[ paramKeys[i] ] );
      }

      if ( innerHTML )
        el.innerHTML = innerHTML;

      return el;
    }

    fillSelect( el, options ) {
      for ( var i=0; i<options.length; i++)
        el.append( this.createElement( "option", { 'value': options[i].id }, options[i].label ) );
    }

  }
  
customElements.define("app-input", AppInputElement);
