
export default () => {
    let _dataImages = {}
    let _mapImg = new Map()
    let _idArray = []
    let _currentImg = {}
    let _image = null
    let _imageOverlay = null

    function resetStrc() {
        _dataImages = {}
        _mapImg = new Map()
        _idArray = []
        _currentImg = {}
        _image = null
        _imageOverlay = null
    }
    return {
        loadData() {
            console.log('load data')
            resetStrc()
            fetch('photos.json', { cache: 'no-store' })
                .then(response => response.json())
                .then((data) => {
                    //console.log('data from fetch: ', data)
                    _dataImages = data;
                    _dataImages = data.images.sort((a, b) => a.id.localeCompare(b.id))
                    let index = 0
                    _dataImages.forEach(item => {
                        _mapImg.set(item.id, { name: item.name, redux: item.redux, caption: item.caption, ix: index })
                        _idArray.push(item.id)
                        index += 1
                    })
                    //console.log('dataimages: ', _dataImages)
                    //console.log('mapImg: ', _mapImg)
                    _image = document.querySelector('#the-image');
                    _imageOverlay = document.querySelector('#image-wrapper');
                    console.log('images data for gallery ok')
                })
                .catch(err => {
                    console.error('error on fetch: ', err)
                });
        },
        displayImage(id) {
            console.log('display image id ', id)
            _currentImg = _mapImg.get(id)
            _image.classList.add('hidden')
            _image.onload = () => { _image.classList.remove('hidden'); }; 
            _image.src = _currentImg.name
            _image.alt = _currentImg.caption
            //console.log('current image ', _currentImg) 
            _imageOverlay.classList.remove('gone');   
        },
        hideGalleryImage(){
            console.log('hide gallery image')
            _imageOverlay.classList.add('gone');
        }
    }
}