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


def generate_template_parser(start_time, forecast_time):
    start_date_time = datetime.datetime.strptime(start_time, "%Y%m%d%H")
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


@click.command()
@click.option("--config-dir", help="config dir")
@click.option("--data-type", help="data type", required=True)
@click.argument("start-time")
@click.argument("forecast-time")
def cli(data_type, config_dir, start_time, forecast_time):
    if config_dir is None:
        config_dir = Path(Path(__file__).parent, "conf")

    config_file_path = find_config(config_dir, data_type)
    if config_file_path is None:
        click.echo('model data type config not found.', err=True)
        click.get_current_context().exit(10)

    config = load_config(config_file_path)
    file_path = find_file(config, start_time, forecast_time)
    click.echo(file_path)


if __name__ == "__main__":
    cli()
