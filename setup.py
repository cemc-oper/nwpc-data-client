from setuptools import setup, find_packages
import codecs
from os import path

here = path.abspath(path.dirname(__file__))

with codecs.open(path.join(here, 'README.md'), encoding='utf-8') as f:
    long_description = f.read()

setup(
    name='nwpc-data-client',

    version='0.1.1',

    description='NWPC data Client',
    long_description=long_description,
    long_description_content_type='text/markdown',

    url='https://github.com/perillaroc/nwpc-operation-system-tool',

    author='perillaroc',
    author_email='perillaroc@gmail.com',

    license='GPL-3.0',

    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: GNU General Public License v3 (GPLv3)',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7'
    ],

    keywords='nwpc data',

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
