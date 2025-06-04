#!/usr/bin/env ruby

require 'bundler/inline'

gemfile do
  source 'https://rubygems.org'
  gem 'webrick'
  gem 'launchy'
  gem 'erb'
end

require 'yaml'
require 'optparse'
require 'webrick'
require 'launchy'
require 'erb'
require 'json'

class SiteDog
  EXAMPLE_CONFIG = <<~YAML
    # Describe your project with a free key-value format, think simple.
    #
    # Random sample:

    registrar: gandi # registrar service
    dns: Route 53 # dns service
    hosting: https://carrd.com # hosting service
    mail: zoho # mail service
  YAML

  TEMPLATE_PATH = 'demo.html.erb'
  DEFAULT_PORT = 8081
  DEFAULT_CONFIG_PATH = './sitedog.yml'

  def self.run
    command = ARGV[0]

    case command
    when 'init'
      init
    when 'demo'
      demo
    when 'help', nil
      show_help
    else
      puts "Unknown command: #{command}"
      show_help
      exit 1
    end
  end

  def self.init
    config_path = DEFAULT_CONFIG_PATH

    OptionParser.new do |opts|
      opts.banner = "Usage: sitedog init [options]"
      opts.on("--config PATH", "Path to config file (default: #{DEFAULT_CONFIG_PATH})") do |path|
        config_path = path
      end
    end.parse!

    if File.exist?(config_path)
      puts "Error: #{config_path} already exists"
      exit 1
    end

    File.write(config_path, EXAMPLE_CONFIG)
    puts "Created #{config_path} configuration file"
  end

  def self.demo
    config_path = DEFAULT_CONFIG_PATH
    port = DEFAULT_PORT

    OptionParser.new do |opts|
      opts.banner = "Usage: sitedog demo [options]"
      opts.on("--config PATH", "Path to config file (default: #{DEFAULT_CONFIG_PATH})") do |path|
        config_path = path
      end
      opts.on("--port PORT", Integer, "Port to run server on (default: #{DEFAULT_PORT})") do |p|
        port = p
      end
    end.parse!

    unless File.exist?(config_path)
      puts "Error: #{config_path} not found. Run 'sitedog init' first."
      exit 1
    end

    # Настраиваем и запускаем сервер
    server = WEBrick::HTTPServer.new(Port: port, AccessLog: [])
    template = ERB.new(File.read(TEMPLATE_PATH))

    server.mount_proc('/') do |req, res|
      config = YAML.load_file(config_path)
      data = SiteDog::Data.new(config)
      res.body = template.result(data.get_binding)
      res['Content-Type'] = 'text/html'
    end

    server.mount_proc('/config') do |req, res|
      config = YAML.load_file(config_path)
      res.body = config.to_json
      res['Content-Type'] = 'application/json'
    end


    puts "Starting demo server at http://localhost:#{port}"
    puts "Press Ctrl+C to stop"

    # Открываем браузер
    browser_thread = Thread.new { Launchy.open("http://localhost:#{port}") }

    # Очищаем временный файл при завершении
    trap('INT') do
      server.shutdown
      Thread.new{sleep(5); exit(0)}
    end

    server.start
  ensure
    server&.shutdown
    browser_thread&.kill
  end

  def self.show_help
    puts "Usage: sitedog <command>"
    puts "\nCommands:"
    puts "  init    Create sitedog.yml configuration file"
    puts "  demo    Start demo server with temporary page"
    puts "  help    Show this help message"
    puts "\nOptions for init:"
    puts "  --config PATH    Path to config file (default: #{DEFAULT_CONFIG_PATH})"
    puts "\nOptions for demo:"
    puts "  --config PATH    Path to config file (default: #{DEFAULT_CONFIG_PATH})"
    puts "  --port PORT      Port to run server on (default: #{DEFAULT_PORT})"
  end
end

class SiteDog::Data
  def initialize(config)
    @config = config
  end

  def get_binding
    binding
  end
end

begin
  SiteDog.run
rescue SystemExit, Interrupt
  puts "Stopping demo server..."
end
