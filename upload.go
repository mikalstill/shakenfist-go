package client

/****

   def get_existing_locks(self):
        r = self._request_url('GET', '/admin/locks')
        return r.json()

    def ping(self, network_ref, address):
        r = self._request_url(
            'GET', '/networks/' + network_ref + '/ping/' + address)
        return r.json()

    def create_upload(self):
        r = self._request_url('POST', '/upload')
        return r.json()

    def send_upload(self, upload_uuid, data):
        r = self._request_url('POST', '/upload/' + upload_uuid,
                              data=data, request_body_is_binary=True)
        return r.json()

    def truncate_upload(self, upload_uuid, offset):
        r = self._request_url(
            'POST', '/upload/' + upload_uuid + '/truncate/' + str(offset))
        return r.json()

****/
