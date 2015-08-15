'use strict';

module.exports = function(grunt) {

    grunt.initConfig({

        jshint: {
            options: {
                jshintrc: '.jshintrc'
            },
            source: [
                'lib.js'
            ],
        },
        uglify: {
            main: {
                beautify: {
                  width: 80,
                  beautify: true,
                  preserveComments:'some'
                },
                files: {
                    'lib.min.js': 'lib.js',
                },
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-jshint');
    grunt.loadNpmTasks('grunt-contrib-uglify');

    grunt.registerTask('build', ['jshint','uglify']);
    grunt.registerTask('default', ['build']);
};