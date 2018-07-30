from setuptools import setup, find_packages
from codecs import open
from os import path

here = path.abspath(path.dirname(__file__))

with open(path.join(here, 'README.md'), encoding='utf-8') as f:
    long_description = f.read()

setup(
    name='nwpc-data-client',

    version='0.1.0',

    description='NWPC data Client',
    long_description=long_description,

    url='https://github.com/perillaroc/nwpc-operation-system-tool',

    author='perillaroc',
    author_email='perillaroc@gmail.com',

    license='GPL-3.0',

    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 3.6'
    ],

    keywords='nwpc slurm client',

    packages=find_packages(exclude=['contrib', 'docs', 'tests']),

    include_package_data=True,
    package_data={
        'nwpc_data_client': ['conf/*.config']
    },

    install_requires=[
        'pyyaml',
        'click'
    ],

    # extras_require={
    #     'test': ['pytest'],
    # },

    entry_points={
        'console_scripts': [
            'nwpc_find_data_path=nwpc_data_client.find_data_path:cli'
        ]
    }
)
