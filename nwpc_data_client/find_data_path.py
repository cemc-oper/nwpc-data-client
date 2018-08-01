# coding: utf-8
import datetime
from pathlib import Path

import click
import yaml


def find_config(config_dir, data_type):
    config_file_path = Path(config_dir, data_type+".yml")
    if config_file_path.is_file():
        return config_file_path
    else:
        return None


def load_config(config_file_path):
    with open(config_file_path) as config_file:
        config = yaml.load(config_file)
        return config


def generate_template_parser(start_date_time, forecast_time):
    YYYY = start_date_time.strftime("%Y")
    MM = start_date_time.strftime("%m")
    DD = start_date_time.strftime("%d")
    HH = start_date_time.strftime("%H")
    FFF = forecast_time

    start_date_time_4dvar = start_date_time - datetime.timedelta(hours=3)
    YYYY_4DV = start_date_time_4dvar.strftime("%Y")
    MM_4DV = start_date_time_4dvar.strftime("%m")
    DD_4DV = start_date_time_4dvar.strftime("%d")
    HH_4DV = start_date_time_4dvar.strftime("%H")

    local_vars = locals()

    def parse_template(template):
        return template.format(**local_vars)

    return parse_template


def find_file(config, start_time, forecast_time):
    parse_template = generate_template_parser(start_time, forecast_time)
    file_name = parse_template(config['file_name'])
    file_path = config['default']
    paths = config['paths']
    for a_path in paths:
        current_dir_path = parse_template(a_path)
        current_file_path = Path(current_dir_path, file_name)
        if current_file_path.is_file():
            file_path = current_file_path
            break

    return file_path


class StartTimeParamType(click.ParamType):
    name = "start_time"

    def convert(self, value, param, ctx):
        try:
            if len(value) != 10:
                self.fail("length of start_time must be 10".format(value=value))
            start_date_time = datetime.datetime.strptime(value, "%Y%m%d%H")
            return start_date_time
        except ValueError:
            self.fail("{value} is not a valid start_time".format(value=value))


class ForecastTimeParamType(click.ParamType):
    name = "forecast_time"

    def convert(self, value, param, ctx):
        try:
            if len(value) > 3:
                self.fail("length of forecast time must less or equal to 3")
            else:
                int_value = int(value)
                return "{value:03}".format(value=int_value)
        except ValueError:
            self.fail("{value} is not a valid forecast_time".format(value=value))


@click.group(invoke_without_command=True)
@click.pass_context
def cli(ctx):
    if ctx.invoked_subcommand is None:
        click.echo(ctx.get_help())


def get_default_local_config_path():
    return Path(Path(__file__).parent, "conf").absolute()


def show_find_local_file_types(ctx, param, value):
    if not value:
        return
    if 'config_dir' in ctx.params:
        config_dir = ctx.params['config_dir']
        config_dir = Path(config_dir).absolute()
    else:
        config_dir = get_default_local_config_path()
    click.echo("config dir:{config_dir}".format(config_dir=str(config_dir)))
    click.echo("file types:")
    for a_config_file in sorted(config_dir.rglob("*.yml")):
        click.echo(a_config_file.relative_to(config_dir).with_suffix(''))
    ctx.exit(0)


@cli.command('local', short_help="Find local data path.")
@click.option("--config-dir", is_eager=True,
              help="Config dir, default: {file_path}".format(file_path=str(get_default_local_config_path())))
@click.option("--data-type", required=True,
              help="Data type used to locate config file path in config dir.")
@click.option("--show-types", is_flag=True, default=False,
              is_eager=True, callback=show_find_local_file_types,
              help="Show supported data types defined in config dir and exit.")
@click.argument("start_time", metavar='<start_time>', type=StartTimeParamType())
@click.argument("forecast_time", metavar='<forecast_time>', type=ForecastTimeParamType())
def find_local_file(config_dir, show_types, data_type, start_time, forecast_time):
    """Find local data path using config files in config dir.

    start_time: YYYYMMDDHH, such as 2018080100
    
    forecast_time: FFF, such as 000
    """

    if config_dir is None:
        config_dir = get_default_local_config_path()

    config_file_path = find_config(config_dir, data_type)
    if config_file_path is None:
        click.echo('model data type config not found.', err=True)
        click.get_current_context().exit(10)

    config = load_config(config_file_path)
    file_path = find_file(config, start_time, forecast_time)
    click.echo(file_path)


if __name__ == "__main__":
    cli()
