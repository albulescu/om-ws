/**
 * Protocol {{.version}}
 * @author Cosmin Albulescu <cosmin@albulescu.ro>
 * @preserve
 */
(function(){

    'use strict';

    var VERSION  = '{{.version}}';
    var ENDPOINT = 'ws://{{.endpoint}}/';


    var events = {};
    var socket;
    var userDisconnected = false;
    var reconnectInterval;
    var reconnectCount;
    var reconnecting=false;

    function bind(event, fct){
        events = events || {};
        events[event] = events[event]   || [];
        events[event].push(fct);
    }

    function unbind(event, fct){
        events = events || {};
        if( event in events === false  ){ return; }
        events[event].splice(events[event].indexOf(fct), 1);
    }

    function trigger(event /* , args... */){
        events = events || {};
        if( event in events === false  ) { return; }
        for(var i = 0; i < events[event].length; i++){
            /* jshint ignore:start */
            events[event][i].apply(this, Array.prototype.slice.call(arguments, 1));
            /* jshint ignore:end */
        }
    }

    function decodePacket( buffer /*ArrayBuffer*/ ) {

        var readString = function(dv, offset, length) {
            var utf16 = new ArrayBuffer(length * 2);
            var utf16View = new Uint16Array(utf16);
            for (var i = 0; i < length; ++i) {
                utf16View[i] = dv.getUint8(offset + i);
            }
            return String.fromCharCode.apply(null, utf16View);
        };

        var data = new DataView(buffer);

        var om = readString(data, 0, 2);

        if( om !== 'om' ) {
            throw new Error('Invalid packet type');
        }

        var packet = {};

        packet.version = data.getUint8(2);

        if( parseInt(packet.version,10) !== parseInt(VERSION,10) ) {
            throw new Error('Invalid version '+packet.version+' in received packet');
        }

        packet.action  = data.getUint8(3);
        packet.size    = data.getUint32(4);
        packet.data    = {};

        var body = new Uint8Array(buffer.slice(8, buffer.byteLength));

        var setValue = function(key, value) {

            if( key.indexOf('.') === -1 ) {
                packet.data[key]=value;
                return;
            }

            var keys = key.split('.');

            var ref = packet.data;

            for (var i = 0; i < keys.length; i++) {

                if( typeof(ref[keys[i]]) !== 'undefined' &&
                    typeof(ref[keys[i]]) !== 'object' ) {
                    throw new Error('Fail to set '+key+' key. Key ' + keys[i] + ' already present as non object');
                }
                if( typeof(ref[keys[i]]) === 'undefined' ) {
                    ref[keys[i]]={};
                }
                if( i < keys.length - 1 ) {
                    if( typeof(ref[keys[i]]) === 'object' ) {
                        ref=ref[keys[i]];
                    }
                } else {
                    ref[keys[i]] = value;
                }
            }
        };

        var key = '', chunk=[], toggle=true, separator = [0xc0,0x80];

        for (var i = 0; i < body.byteLength-1; i++) {
            if( body[i] === separator[0] &&
                body[i+1] === separator[1]) {
                chunk=String.fromCharCode.apply(null, chunk);
                if( toggle ) {
                    key = chunk;
                } else {
                    setValue(key, chunk);
                }
                chunk=[];
                toggle=!toggle;
            } else if(body[i] !== separator[1]) {
                chunk.push(body[i]);
            }
        }

        return packet;
    }

    function encodePacket(action, body) {

        var buffer = new ArrayBuffer(512);
        var data = new DataView(buffer);

        data.setUint8(0,111);//o
        data.setUint8(1,109);//m
        data.setUint8(2,VERSION);//version
        data.setUint8(3,action);//action
        data.setUint32(4,0);//action

        var position=8;
        var addString = function(string) {
            for (var i = 0; i < string.length; i++) {
                data.setUint8(position, string[i].charCodeAt(0));
                position++;
            }
        };

        var addSeparator = function() {
            data.setUint8(position,0xc0);
            position++;
            data.setUint8(position,0x80);
            position++;
        };

        for(var key in body) {
            addString(key);
            addSeparator();
            addString(body[key]+'');
            addSeparator();
        }

        return data.buffer;
    }

    function send(action /*number*/, data /*object*/) {
        if( !connected() ) {
            throw new Error('Not connected');
        }
        socket.send(encodePacket(action, data));
    }

    function reconnect() {
        reconnectCount=0;
        reconnecting=true;
        reconnectInterval = setInterval(function(){
            connect();
            trigger('reconnect', reconnectCount + 1);
            if( reconnectCount>=2 ) {
                reconnecting=false;
                clearInterval(reconnectInterval);
                trigger('error', 'Fail to reconnect');
            }
            reconnectCount++;
        }, 3000);
    }

    function connect() {

        if( connected() ) {
            return;
        }

        socket = new WebSocket(ENDPOINT);
        socket.binaryType = 'arraybuffer';
        socket.onopen = onSocketOpen;
        socket.onclose = onSocketClose;
        socket.onmessage = onMessageReceived;
    }


    function connected() {
        return socket && socket.readyState===1;
    }

    function disconnect() {
        if( connected ) {
            userDisconnected=true;
            socket.close();
        }
    }

    function onSocketOpen() {
        reconnecting=false;
        userDisconnected=false;
        clearInterval(reconnectInterval);
        trigger('connect');
    }

    function onSocketClose() {
        trigger('disconnect');
        if(!userDisconnected && !reconnecting){
            reconnect();
        }
    }

    function onMessageReceived(event) {
        if( event.data !== '' ) {
                try {
                    trigger('data', decodePacket(event.data));
                }
                catch(e) {
                    trigger('error', 'Fail to decode packet:' + e.message);
                }
            }
    }

    var api = {
        'on': bind,
        'off': unbind,
        'send': send,
        'connect': connect,
        'close': disconnect
    };

    Object.defineProperty(api, 'connected', {
        get: function() {
            return connected();
        }
    });

    window.omsocket = api;

    if( typeof(window.omSocketReady) === 'function') {
        /* global omSocketReady */
        omSocketReady(window.omsocket);
    }
})();
