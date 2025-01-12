
export default {
    _data: {},
    _galleryImages: {},
    loadData(callBack) {
        console.log('load data')
        fetch('photos.json', { cache: 'no-store' })
            .then(response => response.json())
            .then((images) => images.sort((a, b) => a.name.localeCompare(b.name)))
            .then((data) => {
                _galleryImages = data;
                console.log('data are: ', data)
                // this.albums = new Map();
                // data.forEach(element => {
                //     if (!this.albums.has(element.textMeta.Directory)) {
                //         this.albums.set(element.textMeta.Directory, new Array());
                //     }
                //     this.albums.get(element.textMeta.Directory).push(element);
                // });
            }).then(() => callBack(this))
            .catch(err => {
                console.error('error on fetch: ', err)
            });
    },
    displayImage(id) {
        console.log('display image id ', id)
    }

}