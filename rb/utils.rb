require 'optparse'
require 'ostruct'
require 'tmpdir'
require 'fileutils'
require 'git'
require 'json'

# OptParse parses CLI options
class OptParse
  # brujula: https://github.com/nogates/brujula
  # ramlrb:  https://github.com/jpb/raml-rb
  @parsers = %w[brujula raml-rb]

  class << self
    attr_accessor :parsers
  end

  def self.parse(args)
    options = OpenStruct.new
    options.verbose = false
    options.parser = ''

    opt_parser = OptionParser.new do |opts|
      opts.banner = 'Usage: ruby main.rb [options]'

      opts.separator ''
      opts.separator 'Specific options:'

      opts.on('--verbose', 'Print errors traces') do |v|
        options.verbose = v
      end

      opts.on('--parser NAME', parsers, 'Parser to test') do |parser|
        options.parser = parser
      end

      opts.separator ''
      opts.separator 'Common options:'

      opts.on_tail('-h', '--help', 'Show this message') do
        puts opts
        exit
      end
    end

    opt_parser.parse!(args)
    options
  end
end

# Clones raml-tck repo and returns its tests path
def clone_tck_repo
  repo_dir = File.join(Dir.tmpdir, 'raml-tck')
  FileUtils.remove_dir(repo_dir) if File.directory?(repo_dir)
  puts "Cloning raml-tck repo to #{repo_dir}"
  repo = Git.clone(
    'git@github.com:raml-org/raml-tck.git',
    '', path: repo_dir
  )
  repo.checkout('rename-cleanup')
  File.join(repo_dir, 'tests', 'raml-1.0')
end

# Lists RAML files in a folder
def list_ramls(ex_dir)
  manifest_path = File.join(ex_dir, 'manifest.json')
  manifest_file = File.read(manifest_path)
  manifest = JSON.parse(manifest_file)
  manifest['filePaths'].map do |fpath|
    File.join(ex_dir, fpath)
  end
end

# Saves report to JSON file
def save_report(report)
  reports_dir = File.join(__dir__, '..', 'reports', 'json')
  FileUtils.mkdir_p(reports_dir) unless File.directory?(reports_dir)
  report_file = File.join(reports_dir, "#{report['parser']}.json")
  File.open(report_file, 'w') do |f|
    f.write(JSON.pretty_generate(report))
  end
end
