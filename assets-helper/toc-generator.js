import fs from 'fs';
import path from 'path';
import ExifReader from 'exifreader';

const assets_directory = '../assets';
const output_filename = 'files.json';
const skip_file = ['.DS_Store'];

async function getAllFilesRecursively(directory) {
  let files = [];

  try {
    const childFiles = await fs.promises.readdir(directory);

    for (let file of childFiles) {
      const fullPath = path.join(directory, file);
      const relativePath = fullPath.replace(assets_directory, '');

      if (fs.statSync(fullPath).isDirectory()) {
        files = files.concat(await getAllFilesRecursively(fullPath));
      } else if (!(skip_file.includes(file))) {
        files.push(relativePath)
      }
    }
  } catch (err) {
    console.error("Error: ", err);
  }

  return files;
}

async function groupFilesByDirectory(files) {
  const group_actions = files.map(async (val) => {
    const path_split = val.split('/');
    const path = path_split.slice(1).reverse().slice(1);
    const filename = path_split.reverse()[0];

    const basic_info = {
      path: path ? `/${path.reverse().join('/')}` : '/',
      filename,
    };

    try {
      const { exif } = await ExifReader.load(assets_directory + val, {
        expanded: true,
      });

      if (exif) {
        return {
          ...basic_info,
          exif: {
            datetime: exif.DateTime?.value[0] ?? undefined,
            shutter: exif.ExposureTime?.description ?? undefined,
            fstop: exif.FNumber?.description ?? undefined,
            iso: exif.ISOSpeedRatings?.description ?? undefined,
            focal: exif.FocalLengthIn35mmFilm?.description ?? undefined,
            equipment: {
              camera: exif.Make || exif.Model ? `${exif.Make?.description} ${exif.Model?.description}` : undefined,
              lens: exif.LensModel?.description ?? undefined,
            },
          },
        };
      }
    } catch (err) {
      console.error(err);
    }

    return basic_info;

  });

  const groups = await Promise.all(group_actions);

  const group_by = Object.groupBy(groups, ({ path }) => path);

  return formatGroups(group_by);
}

async function formatGroups(groups) {
  for (let key in groups) {
    groups[key] = groups[key].map(({filename, exif}) => {
      return {
        filename,
        exif,
      };
    });
  }

  return groups;
}

async function main() {
  let allFiles = await getAllFilesRecursively(assets_directory);

  let grouped = await groupFilesByDirectory(allFiles);

  fs.writeFile(assets_directory + '/' + output_filename, JSON.stringify(grouped), 'utf8', (err) => {
    if (err) {
      console.log("Error writing file");
      return;
    }
  });
}

main().catch(console.error);
